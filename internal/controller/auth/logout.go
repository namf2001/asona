package auth

import (
	"context"
	"errors"

	"asona/internal/repository/sessions"

	pkgerrors "github.com/pkg/errors"
)

// Logout invalidates an authenticated user's session token.
func (i impl) Logout(ctx context.Context, token string) error {
	err := i.repo.Session().Delete(ctx, token)
	if err != nil {
		if errors.Is(err, sessions.ErrSessionNotFound) {
			return pkgerrors.WithStack(ErrSessionNotFound)
		}
		return err
	}

	return nil
}
