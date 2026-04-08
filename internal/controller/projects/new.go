package projects

import (
	"context"

	"asona/internal/repository"
)

// Controller defines the business logic operations for project management.
type Controller interface {
	// Create creates a new project within a workplace.
	Create(ctx context.Context, userID int64, input CreateProjectInput) (ProjectResponse, error)
	// GetByID retrieves a single project by its unique identifier.
	GetByID(ctx context.Context, projectID int64) (ProjectResponse, error)
	// ListByWorkplace retrieves all projects belonging to a specific workplace.
	ListByWorkplace(ctx context.Context, workplaceID int64) ([]ProjectResponse, error)
}

type impl struct {
	repo repository.Registry
}

// New returns a new project Controller instance.
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
