package accounts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

// GetByProvider retrieves an OAuth account by provider and provider account ID.
func (i impl) GetByProvider(ctx context.Context, provider, providerAccountID string) (model.Account, error) {
	var account model.Account
	var tokenExpiresAt sql.NullTime

	query := `
		SELECT id, user_id, provider, provider_account_id, access_token, refresh_token,
		       token_expires_at, id_token, scope, created_at, updated_at
		FROM public.auth_providers
		WHERE provider = $1 AND provider_account_id = $2
	`

	err := i.db.QueryRowContext(ctx, query, provider, providerAccountID).Scan(
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderAccountID,
		&account.AccessToken,
		&account.RefreshToken,
		&tokenExpiresAt,
		&account.IDToken,
		&account.Scope,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Account{}, pkgerrors.WithStack(ErrAccountNotFound)
		}
		return model.Account{}, pkgerrors.WithStack(fmt.Errorf("failed to get oauth account: %w", err))
	}

	if tokenExpiresAt.Valid {
		account.TokenExpiresAt = &tokenExpiresAt.Time
	}

	return account, nil
}
