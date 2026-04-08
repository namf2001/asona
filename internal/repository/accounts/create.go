package accounts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

// Create links an OAuth account to a user.
func (i impl) Create(ctx context.Context, account model.Account) (model.Account, error) {
	if account.CreatedAt.IsZero() {
		account.CreatedAt = time.Now()
	}
	if account.UpdatedAt.IsZero() {
		account.UpdatedAt = time.Now()
	}

	var tokenExpiresAt any = nil
	if account.TokenExpiresAt != nil {
		tokenExpiresAt = *account.TokenExpiresAt
	}

	query := `
		INSERT INTO public.auth_providers (
			user_id, provider, provider_account_id, access_token, refresh_token,
			token_expires_at, id_token, scope, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := i.db.QueryRowContext(ctx, query,
		account.UserID,
		account.Provider,
		account.ProviderAccountID,
		account.AccessToken,
		account.RefreshToken,
		tokenExpiresAt,
		account.IDToken,
		account.Scope,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&account.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return model.Account{}, pkgerrors.WithStack(ErrAccountAlreadyExists)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return model.Account{}, pkgerrors.WithStack(ErrAccountNotFound)
		}
		return model.Account{}, pkgerrors.WithStack(fmt.Errorf("failed to create oauth account: %w", err))
	}

	return account, nil
}
