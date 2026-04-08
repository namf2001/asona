package projects

import (
	"context"

	"asona/internal/model"
)

// GetByID retrieves a specific project by its ID.
func (i impl) GetByID(ctx context.Context, id int64) (model.Project, error) {
	return model.Project{}, nil
}
