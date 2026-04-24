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
	ctrlChannel "asona/internal/controller/chat/channel"
	ctrlMessage "asona/internal/controller/chat/message"
	"asona/internal/controller/organizations"
	"asona/internal/controller/projects"
	"asona/internal/controller/tasks"
	"asona/internal/controller/workplaces"
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
	defer func() { _ = logger.Sync() }()

	// Initialize controllers.
	// Controllers are passed directly into the router; handlers are constructed
	// inline inside each route-group function (following the Thor pattern).
	authCtrl := auth.New(repo, s.oauth, s.mail)
	orgCtrl := organizations.New(repo)
	channelCtrl := ctrlChannel.New(repo)
	messageCtrl := ctrlMessage.New(repo)
	projCtrl := projects.New(repo)
	taskCtrl := tasks.New(repo)
	wpCtrl := workplaces.New(repo)

	rtr := router{
		ctx:         s.ctx,
		db:          s.db,
		rdb:         s.rdb,
		logger:      logger,
		ws:          s.ws,
		authCtrl:    authCtrl,
		orgCtrl:     orgCtrl,
		channelCtrl: channelCtrl,
		messageCtrl: messageCtrl,
		projCtrl:    projCtrl,
		taskCtrl:    taskCtrl,
		wpCtrl:      wpCtrl,
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
