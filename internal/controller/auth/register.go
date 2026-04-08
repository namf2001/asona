package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"asona/internal/model"
	"asona/internal/repository/users"

	pkgerrors "github.com/pkg/errors"
)

// RegisterInput represents the input data for new user registration.
type RegisterInput struct {
	Name     string `json:"name"     binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register creates a new user account in the system.
func (i impl) Register(ctx context.Context, input RegisterInput) (model.User, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, pkgerrors.WithStack(fmt.Errorf("failed to hash password: %w", err))
	}

	user := model.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPasswordBytes),
	}

	createdUser, err := i.repo.User().Create(ctx, user)
	if err != nil {
		if errors.Is(err, users.ErrEmailAlreadyExists) {
			return model.User{}, pkgerrors.WithStack(ErrUserAlreadyExists)
		}
		return model.User{}, pkgerrors.WithStack(fmt.Errorf("failed to register user: %w", err))
	}

	return createdUser, nil
}
