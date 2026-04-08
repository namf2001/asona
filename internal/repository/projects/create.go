package projects

import (
	"context"

	"asona/internal/model"
)

// Create inserts a new project into the database.
func (i impl) Create(ctx context.Context, project model.Project) (model.Project, error) {
	return project, nil
}
