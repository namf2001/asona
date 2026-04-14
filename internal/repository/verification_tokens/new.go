package verification_tokens

import (
	"context"

	"asona/internal/model"
	"asona/internal/repository/db/pg"
)

type Repository interface {
	Create(ctx context.Context, token model.VerificationToken) error
	GetValidToken(ctx context.Context, identifier, tokenStr string) (model.VerificationToken, error)
	Delete(ctx context.Context, identifier, tokenStr string) error
	DeleteAllForIdentifier(ctx context.Context, identifier string) error
}

type impl struct {
	db pg.ContextExecutor
}

// New creates a new verification_tokens repository
func New(db pg.ContextExecutor) Repository {
	return impl{db: db}
}
