package auth

import "asona/internal/controller/auth"

// Handler handles HTTP requests for authentication.
type Handler struct {
	ctrl auth.Controller
}

// New creates a new authentication handler.
func New(ctrl auth.Controller) Handler {
	return Handler{
		ctrl: ctrl,
	}
}
