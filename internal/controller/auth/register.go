package auth

import (
	"context"
	"fmt"
	"asona/internal/model"
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
	// 1. Hash password (omitted for brevity, assume bcrypt)
	hashedPassword := input.Password 

	user := model.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// 2. Map to repository
	createdUser, err := i.repo.User().Create(ctx, user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to register user: %w", err)
	}

	return createdUser, nil
}
