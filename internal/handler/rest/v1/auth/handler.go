package auth

import (
	"asona/internal/controller/auth"
	"asona/internal/service/redis"
)

// Handler handles HTTP requests for authentication.
type Handler struct {
	ctrl auth.Controller
	rdb  redis.Service
}

// New creates a new authentication handler.
func New(ctrl auth.Controller, rdb redis.Service) Handler {
	return Handler{
		ctrl: ctrl,
		rdb:  rdb,
	}
}
