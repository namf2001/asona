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
}

type impl struct {
	db pg.ContextExecutor
}

// New returns a new task Repository.
func New(db pg.ContextExecutor) Repository {
	return impl{db: db}
}
