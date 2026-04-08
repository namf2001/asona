package tasks

import (
	"context"

	"asona/internal/model"
)

// Create inserts a new task into the database.
func (i impl) Create(ctx context.Context, task model.Task) (model.Task, error) {
	return task, nil
}
