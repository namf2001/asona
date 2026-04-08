package projects

import (
	"asona/internal/controller/projects"
)

// Handler handles HTTP requests related to project management.
type Handler struct {
	ctrl projects.Controller
}

// New creates a new project handler instance.
func New(ctrl projects.Controller) Handler {
	return Handler{ctrl: ctrl}
}
