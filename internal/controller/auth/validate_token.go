package auth

import (
	"context"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/pkg/jwt"
	"asona/internal/repository/sessions"
)

// ValidateToken verifies the JWT signature and confirms the backing session
// still exists in the database. This is the single point of token validation
// used by the auth middleware — both expiry and revocation are enforced here.
func (i impl) ValidateToken(ctx context.Context, token string) (int64, string, error) {
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return 0, "", pkgerrors.WithStack(fmt.Errorf("%w: %v", ErrInvalidToken, err))
	}

	// Confirm the session is still present and unexpired in the database.
	// GetByToken already filters WHERE expires_at > NOW(), so a deleted (logged-out)
	// or naturally expired session will return ErrSessionNotFound here.
	if _, err := i.repo.Session().GetByToken(ctx, token); err != nil {
		if errors.Is(err, sessions.ErrSessionNotFound) {
			return 0, "", pkgerrors.WithStack(ErrSessionNotFound)
		}
		return 0, "", err
	}

	return claims.UserID, claims.Email, nil
}
