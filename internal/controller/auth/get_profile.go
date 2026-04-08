package auth

import (
	"context"

	"asona/internal/model"
)

// GetProfile retrieves a user's profile information by their ID.
func (i impl) GetProfile(ctx context.Context, userID int64) (model.User, error) {
	return i.repo.User().GetByID(ctx, userID)
}
