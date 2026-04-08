package workplaces

import (
	"asona/internal/model"
	"asona/internal/repository/db/pg"
	"context"
)

// Repository defines workplace data operations.
type Repository interface {
	// Create creates a new workplace within an organization
	Create(ctx context.Context, wp model.Workplace) (model.Workplace, error)

	// GetByID retrieves a workplace by ID
	GetByID(ctx context.Context, id int64) (model.Workplace, error)

	// List retrieves workplaces with optional filters
	List(ctx context.Context, filter ListFilter) ([]model.Workplace, error)

	// AddMember adds a user to a workplace
	AddMember(ctx context.Context, member model.WorkplaceMember) error

	// RemoveMember removes a user from a workplace
	RemoveMember(ctx context.Context, workplaceID, userID int64) error
}

// ListFilter provides criteria for filtering workplace lists.
type ListFilter struct {
	CreatedBy int64
}

type impl struct {
	db pg.ContextExecutor
}

// New returns a new workplace Repository instance.
func New(db pg.ContextExecutor) Repository {
	return impl{
		db: db,
	}
}
