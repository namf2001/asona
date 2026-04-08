package organizations

import (
	"asona/internal/controller/organizations"
)

// Handler handles HTTP requests for organization management.
type Handler struct {
	ctrl organizations.Controller
}

// New creates a new organization handler.
func New(ctrl organizations.Controller) Handler {
	return Handler{
		ctrl: ctrl,
	}
}
