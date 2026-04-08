package accounts

import (
	"asona/internal/model"
	"asona/internal/repository/db/pg"
	"context"
)

type Repository interface {
	// Create links an OAuth account to a user
	Create(ctx context.Context, account model.Account) (model.Account, error)

	// GetByProvider retrieves an account by provider and provider ID
	GetByProvider(ctx context.Context, provider, providerAccountID string) (model.Account, error)
}

type impl struct {
	db pg.ContextExecutor
}

func New(db pg.ContextExecutor) Repository {
	return impl{
		db: db,
	}
}
