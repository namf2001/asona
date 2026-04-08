package auth

import (
	"context"
	"errors"

	"asona/internal/model"
	"asona/internal/repository/users"

	pkgerrors "github.com/pkg/errors"
)

// GetProfile retrieves a user's profile information by their ID.
func (i impl) GetProfile(ctx context.Context, userID int64) (model.User, error) {
	user, err := i.repo.User().GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return model.User{}, pkgerrors.WithStack(ErrUserNotFound)
		}
		return model.User{}, err
	}

	return user, nil
}
