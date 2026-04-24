package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	_ "asona/docs/swagger"

	"github.com/gin-contrib/cors"
	ginsessions "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	authctrl "asona/internal/controller/auth"
	ctrlChannel "asona/internal/controller/chat/channel"
	ctrlMessage "asona/internal/controller/chat/message"
	"asona/internal/controller/organizations"
	"asona/internal/controller/projects"
	"asona/internal/controller/tasks"
	"asona/internal/controller/workplaces"
	"asona/internal/handler/middleware"
	handlerauth "asona/internal/handler/rest/v1/auth"
	handlerchat "asona/internal/handler/rest/v1/chat"
	handlerorg "asona/internal/handler/rest/v1/organizations"
	handlerproj "asona/internal/handler/rest/v1/projects"
	handlertask "asona/internal/handler/rest/v1/tasks"
	handlerws "asona/internal/handler/rest/v1/websocket"
	handlerwp "asona/internal/handler/rest/v1/workplaces"
	"asona/internal/service/database"
	"asona/internal/service/redis"
	wsservice "asona/internal/service/websocket"

	"go.uber.org/zap"
)

// router defines the routes & handlers of the app.
type router struct {
	ctx    context.Context
	db     database.Service
	rdb    redis.Service
	logger *zap.Logger
	ws     wsservice.Service

	authCtrl    authctrl.Controller
	orgCtrl     organizations.Controller
	channelCtrl ctrlChannel.Controller
	messageCtrl ctrlMessage.Controller
	projCtrl    projects.Controller
	taskCtrl    tasks.Controller
	wpCtrl      workplaces.Controller
}

// handler returns the http.Handler for use by the server.
func (rtr router) handler() http.Handler {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(rtr.logger))
	r.Use(gin.Recovery())
	r.Use(ginsessions.Sessions("session", cookie.NewStore([]byte(createRandStr()))))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	rtr.routes(r)

	return r
}

func (rtr router) routes(r *gin.Engine) {
	rtr.public(r)
	rtr.authenticated(r)
	rtr.websocketRoutes(r)
}

func (rtr router) public(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "Hello World"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, rtr.db.Health())
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authHandler := handlerauth.New(rtr.authCtrl, rtr.rdb)

	v1 := r.Group("/api/v1")
	v1.POST("/register", middleware.RSAAuthMiddleware(), authHandler.Register)
	v1.POST("/login", middleware.RSAAuthMiddleware(), authHandler.Login)
	v1.GET("/auth/google", authHandler.GoogleLogin)
	v1.GET("/auth/google/callback", authHandler.GoogleCallback)
	// ExchangeOAuthCode exchanges a one-time Redis code for the real JWT session token.
	// Token is never exposed in a URL — only returned via JSON body.
	v1.POST("/auth/exchange", authHandler.ExchangeOAuthCode)
}

func (rtr router) authenticated(r *gin.Engine) {
	authHandler := handlerauth.New(rtr.authCtrl, rtr.rdb)
	orgHandler := handlerorg.New(rtr.orgCtrl)
	chatHandler := handlerchat.New(rtr.channelCtrl, rtr.messageCtrl)
	projHandler := handlerproj.New(rtr.projCtrl)
	taskHandler := handlertask.New(rtr.taskCtrl)
	wpHandler := handlerwp.New(rtr.wpCtrl)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.TokenCheckMiddleware(rtr.authCtrl))
	v1.Use(middleware.RSAAuthMiddleware())

	// Auth
	v1.GET("/profile", authHandler.Profile)
	v1.GET("/me/onboarding", authHandler.GetOnboardingState)
	v1.POST("/logout", authHandler.Logout)
	v1.PATCH("/me/onboard", authHandler.CompleteOnboard)

	// Organizations
	orgs := v1.Group("/organizations")
	{
		orgs.POST("", orgHandler.CreateOrganization)
		orgs.GET("/:id", orgHandler.GetOrganization)
	}

	// Workplaces
	wps := v1.Group("/workplaces")
	{
		wps.POST("", wpHandler.CreateWorkplace)
		wps.GET("/:id/projects", projHandler.ListProjects)
	}

	// Channels
	channels := v1.Group("/channels")
	{
		channels.POST("", chatHandler.CreateChannel)
		channels.GET("/:id", chatHandler.GetChannel)
		channels.GET("/:id/messages", chatHandler.ListMessages)
	}

	// Projects
	projs := v1.Group("/projects")
	{
		projs.POST("", projHandler.CreateProject)
		projs.GET("/:id", projHandler.GetProject)
		projs.GET("/:id/tasks", taskHandler.ListTasks)
	}

	// Tasks
	taskGroup := v1.Group("/tasks")
	{
		taskGroup.POST("", taskHandler.CreateTask)
		taskGroup.GET("/:id", taskHandler.GetTask)
		taskGroup.PUT("/:id", taskHandler.UpdateTask)
	}

	// Messages
	v1.POST("/messages", chatHandler.SendMessage)
}

// websocketRoutes registers WebSocket routes.
func (rtr router) websocketRoutes(r *gin.Engine) {
	wsHandler := handlerws.New(rtr.ws)

	// Public WebSocket (userId via query param)
	r.GET("/ws", wsHandler.HandleWebSocket)

	// Authenticated WebSocket
	wsAuth := r.Group("/api/v1/ws")
	wsAuth.Use(middleware.TokenCheckMiddleware(rtr.authCtrl))
	wsAuth.GET("", wsHandler.HandleWebSocket)
}

// createRandStr generates a random 32-byte hex string for use as a session secret.
func createRandStr() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("server: failed to generate random session key: " + err.Error())
	}
	return hex.EncodeToString(b)
}
