package sessions

import (
	"asona/internal/model"
	"asona/internal/repository/db/pg"
	"context"
)

type Repository interface {
	// Create creates a new session
	Create(ctx context.Context, session model.Session) (model.Session, error)

	// GetByToken retrieves a session by its token
	GetByToken(ctx context.Context, token string) (model.Session, error)

	// Delete deletes a session by token
	Delete(ctx context.Context, token string) error
}

type impl struct {
	db pg.ContextExecutor
}

func New(db pg.ContextExecutor) Repository {
	return impl{
		db: db,
	}
}
