package projects

import (
	"context"

	"asona/internal/model"
)

// List fetches projects from the database based on the provided filter criteria.
func (i impl) List(ctx context.Context, filter ListFilter) ([]model.Project, error) {
	// Implementation will be added later during DB logic phase
	return nil, nil
}
