package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"asona/config"
	"asona/internal/controller/auth"
	handlerauth "asona/internal/handler/rest/v1/auth"
	handlerws "asona/internal/handler/rest/v1/websocket"
	"asona/internal/repository"
	"asona/internal/service/database"
	"asona/internal/service/mail"
	"asona/internal/service/oauth"
	"asona/internal/service/redis"
	"asona/internal/service/s3"
	"asona/internal/service/websocket"
)

// Server holds the application's HTTP server and service dependencies.
type Server struct {
	ctx   context.Context
	port  int
	db    database.Service
	rdb   redis.Service
	mail  mail.Service
	oauth oauth.Service
	s3    s3.Service
	ws    websocket.Service
}

// NewServer creates a new HTTP server with the given configuration.
func NewServer() *http.Server {
	cfg := config.GetConfig()
	port, _ := strconv.Atoi(cfg.AppPort)

	ctx := context.Background()
	rdb := redis.New()

	s := &Server{
		ctx:   ctx,
		port:  port,
		db:    database.New(),
		rdb:   rdb,
		mail:  mail.New(),
		oauth: oauth.New(),
		s3:    s3.New(),
		ws:    websocket.New(ctx, rdb.WSClient()),
	}

	// Start WebSocket service
	s.ws.Run()

	// Initialize repositories
	repo := repository.New(s.db.DB())

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize handlers
	authCtrl := auth.New(repo)
	authHandler := handlerauth.New(authCtrl)
	wsHandler := handlerws.New(s.ws)

	rtr := router{
		ctx:         s.ctx,
		db:          s.db,
		rdb:         s.rdb,
		logger:      logger,
		authHandler: authHandler,
		wsHandler:   wsHandler,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      rtr.handler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
