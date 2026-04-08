package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	_ "asona/docs/swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"asona/internal/handler/middleware"
	"asona/internal/handler/rest/v1/auth"
	"asona/internal/handler/rest/v1/chat"
	"asona/internal/handler/rest/v1/organizations"
	"asona/internal/handler/rest/v1/websocket"
	"asona/internal/service/database"
	"asona/internal/service/redis"

	"go.uber.org/zap"
)

// router defines the routes & handlers of the app.
type router struct {
	ctx    context.Context
	db     database.Service
	rdb    redis.Service
	logger *zap.Logger

	authHandler auth.Handler
	orgHandler  organizations.Handler
	chatHandler chat.Handler
	wsHandler   *websocket.Handler
}

// handler returns the http.Handler for use by the server.
func (rtr router) handler() http.Handler {
	r := gin.New()

	// Middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(rtr.logger))
	r.Use(gin.Recovery())

	r.Use(sessions.Sessions("session", cookie.NewStore([]byte(createRandStr()))))

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

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// api/v1 public routes
	v1 := r.Group("/api/v1")
	v1.POST("/login", middleware.RSAAuthMiddleware(), rtr.authHandler.Login)
}

func (rtr router) authenticated(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.Use(middleware.TokenCheckMiddleware())
	v1.Use(middleware.RSAAuthMiddleware())

	v1.GET("/profile", rtr.authHandler.Profile)

	// Organizations
	orgs := v1.Group("/organizations")
	{
		orgs.POST("", rtr.orgHandler.Create)
		orgs.GET("/:id", rtr.orgHandler.Get)
	}

	// Channels
	channels := v1.Group("/channels")
	{
		channels.POST("", rtr.chatHandler.CreateChannel)
		channels.GET("/:id", rtr.chatHandler.GetChannel)
		channels.GET("/:id/messages", rtr.chatHandler.ListMessages)
	}

	// Messages
	v1.POST("/messages", rtr.chatHandler.SendMessage)
}

// websocketRoutes registers WebSocket routes.
func (rtr router) websocketRoutes(r *gin.Engine) {
	// WebSocket endpoint (public, requires userId query param)
	r.GET("/ws", rtr.wsHandler.HandleWebSocket)

	// WebSocket endpoint with authentication
	wsAuth := r.Group("/api/v1/ws")
	wsAuth.Use(middleware.TokenCheckMiddleware())
	wsAuth.GET("", rtr.wsHandler.HandleWebSocket)
}

// createRandStr generates a random 32-byte hex string for use as a session secret.
func createRandStr() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("server: failed to generate random session key: " + err.Error())
	}
	return hex.EncodeToString(b)
}
