package auth

import (
	"context"
	"errors"

	"asona/internal/model"
	"asona/internal/repository/users"
	pkgerrors "github.com/pkg/errors"
)

// LoginInput represents the input data for user authentication.
type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user and returns their information and a session token.
func (i impl) Login(ctx context.Context, input LoginInput) (model.User, string, error) {
	// 1. Get user by email
	user, err := i.repo.User().GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return model.User{}, "", pkgerrors.WithStack(ErrUserNotFound)
		}
		return model.User{}, "", err
	}

	// 2. Mock password check
	// 3. Mock session token generation
	sessionToken := "mock-session-token-" + input.Email

	return user, sessionToken, nil
}
