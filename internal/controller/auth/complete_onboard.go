package auth

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/repository/users"
)

// CompleteOnboard marks the authenticated user as having finished the onboarding flow.
func (i impl) CompleteOnboard(ctx context.Context, userID int64) error {
	if err := i.repo.User().CompleteOnboard(ctx, userID); err != nil {
		if pkgerrors.Is(err, users.ErrUserNotFound) {
			return pkgerrors.WithStack(ErrUserNotFound)
		}
		return err
	}
	return nil
}
