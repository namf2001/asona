package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"context"
	"asona/config"
	"asona/internal/controller/auth"
	handlerauth "asona/internal/handler/rest/v1/auth"
	"asona/internal/service/database"
	"asona/internal/service/mail"
	"asona/internal/service/oauth"
	"asona/internal/service/redis"
	"asona/internal/service/s3"
	"asona/internal/repository"
	"go.uber.org/zap"
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
}

// NewServer creates a new HTTP server with the given configuration.
func NewServer() *http.Server {
	cfg := config.GetConfig()
	port, _ := strconv.Atoi(cfg.AppPort)
	
	s := &Server{
		ctx:  context.Background(),
		port: port,

		db:    database.New(),
		rdb:   redis.New(),
		mail:  mail.New(),
		oauth: oauth.New(),
		s3:    s3.New(),
	}

	// Initialize repositories
	repo := repository.New(s.db.DB())

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize handlers
	authCtrl := auth.New(repo)
	authHandler := handlerauth.New(authCtrl)

	rtr := router{
		ctx:             s.ctx,
		db:              s.db,
		rdb:             s.rdb,
		logger:          logger,
		
		authHandler:     authHandler,
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
