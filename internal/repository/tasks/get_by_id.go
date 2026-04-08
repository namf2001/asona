package tasks

import (
	"context"

	"asona/internal/model"
)

// GetByID retrieves a specific task by its ID.
func (i impl) GetByID(ctx context.Context, id int64) (model.Task, error) {
	return model.Task{}, nil
}
