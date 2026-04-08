package auth

import (
	"asona/internal/model"
	"asona/internal/repository/users"
	"context"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// LoginInput represents the input data for user authentication.
type LoginInput struct {
	Email     string `json:"email"    binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	UserAgent string
	IPAddress string
}

// Login authenticates a user and returns their information and a session token.
func (i impl) Login(ctx context.Context, input LoginInput) (model.User, string, error) {
	user, err := i.repo.User().GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return model.User{}, "", pkgerrors.WithStack(ErrUserNotFound)
		}
		return model.User{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return model.User{}, "", pkgerrors.WithStack(ErrInvalidPassword)
		}
		return model.User{}, "", pkgerrors.WithStack(fmt.Errorf("failed to compare password: %w", err))
	}

	sessionToken, err := i.issueSession(ctx, user, input.UserAgent, input.IPAddress)
	if err != nil {
		return model.User{}, "", err
	}

	return user, sessionToken, nil
}
