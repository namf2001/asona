package tasks

import (
	"context"

	"asona/internal/repository"
)

// Controller defines the business logic operations for task management.
type Controller interface {
	// Create creates a new task within a project.
	Create(ctx context.Context, senderID int64, input CreateTaskInput) (TaskResponse, error)
	// GetByID retrieves a single task by its unique identifier.
	GetByID(ctx context.Context, taskID int64) (TaskResponse, error)
	// Update updates an existing task's details.
	Update(ctx context.Context, senderID int64, taskID int64, input UpdateTaskInput) error
	// ListByProject retrieves all tasks belonging to a specific project.
	ListByProject(ctx context.Context, projectID int64) ([]TaskResponse, error)
}

type impl struct {
	repo repository.Registry
}

// New returns a new task Controller instance.
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
