package projects

import (
	"context"

	"asona/internal/model"
	"asona/internal/repository/db/pg"
)

// Repository defines the interface for project data operations.
type Repository interface {
	// Create inserts a new project into the database.
	Create(ctx context.Context, project model.Project) (model.Project, error)
	// GetByID retrieves a specific project by its ID.
	GetByID(ctx context.Context, id int64) (model.Project, error)
	// AddMember adds a existing user as a member of a project.
	AddMember(ctx context.Context, member model.ProjectMember) error
}

type impl struct {
	db pg.ContextExecutor
}

// New returns a new project Repository.
func New(db pg.ContextExecutor) Repository {
	return impl{db: db}
}
