package users

import (
	"asona/internal/model"
	"asona/internal/repository/db/pg"
	"context"
)

type ListFilters struct {
	Limit  int
	Offset int
}

type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user model.User) (model.User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (model.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (model.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user model.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id int64) error

	// CompleteOnboard marks the user as having completed the onboarding flow
	// by setting onboarded_at to the current timestamp.
	CompleteOnboard(ctx context.Context, userID int64) error
}

type impl struct {
	db pg.ContextExecutor
}

func New(db pg.ContextExecutor) Repository {
	return impl{
		db: db,
	}
}
