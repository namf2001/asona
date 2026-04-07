package server

import (
	"context"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "asona/docs/swagger"

	"asona/internal/handler/middleware"
	"asona/internal/handler/rest/v1/auth"
	"asona/internal/service/database"
	"asona/internal/service/redis"
	"go.uber.org/zap"
)

// router defines the routes & handlers of the app.
type router struct {
	ctx context.Context
	db  database.Service
	rdb redis.Service
	logger *zap.Logger

	authHandler     *auth.Handler
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
}

// createRandStr generates a random 32-byte hex string for use as a session secret.
func createRandStr() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("server: failed to generate random session key: " + err.Error())
	}
	return hex.EncodeToString(b)
}
