package workplaces

import (
	"asona/internal/controller/workplaces"
)

// Handler handles HTTP requests for workplace management.
type Handler struct {
	ctrl workplaces.Controller
}

// New creates a new workplace handler.
func New(ctrl workplaces.Controller) Handler {
	return Handler{
		ctrl: ctrl,
	}
}
