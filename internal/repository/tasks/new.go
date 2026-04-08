package tasks

import (
	"context"

	"asona/internal/model"
	"asona/internal/repository/db/pg"
)

// Repository defines the interface for task data operations.
type Repository interface {
	// Create inserts a new task into the database.
	Create(ctx context.Context, task model.Task) (model.Task, error)
	// GetByID retrieves a specific task by its ID.
	GetByID(ctx context.Context, id int64) (model.Task, error)
	// Update updates an existing task's information.
	Update(ctx context.Context, task model.Task) error
	// List retrieves tasks with filtering.
	List(ctx context.Context, filter ListFilter) ([]model.Task, error)
}

// ListFilter represents criteria for filtering tasks.
type ListFilter struct {
	ProjectID int64
	Assignee  int64
	Status    string
}

type impl struct {
	db pg.ContextExecutor
}

// New returns a new task Repository.
func New(db pg.ContextExecutor) Repository {
	return impl{db: db}
}
