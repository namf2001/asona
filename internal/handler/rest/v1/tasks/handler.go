package tasks

import (
	"asona/internal/controller/tasks"
)

// Handler handles HTTP requests related to task management.
type Handler struct {
	ctrl tasks.Controller
}

// New creates a new task handler instance.
func New(ctrl tasks.Controller) Handler {
	return Handler{ctrl: ctrl}
}
