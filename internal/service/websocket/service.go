package websocket

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Service defines WebSocket operations.
type Service interface {
	// Hub returns the WebSocket hub.
	Hub() *Hub

	// Upgrader returns the WebSocket upgrader.
	Upgrader() *websocket.Upgrader

	// Run starts the WebSocket service.
	Run()
}

type service struct {
	ctx      context.Context
	hub      *Hub
	upgrader *websocket.Upgrader
}

// New creates a new WebSocket Service.
func New(ctx context.Context, redisClient *redis.Client) Service {
	hub := NewHub(ctx, redisClient)

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Add proper origin check in production
			return true
		},
	}

	return service{
		ctx:      ctx,
		hub:      hub,
		upgrader: upgrader,
	}
}

// Hub returns the WebSocket hub.
func (s service) Hub() *Hub {
	return s.hub
}

// Upgrader returns the WebSocket upgrader.
func (s service) Upgrader() *websocket.Upgrader {
	return s.upgrader
}

// Run starts the WebSocket service.
func (s service) Run() {
	go s.hub.Run()
}
