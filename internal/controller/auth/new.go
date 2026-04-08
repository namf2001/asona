package auth

import (
	"asona/internal/model"
	"asona/internal/repository"
	"context"
)

// Controller defines auth business logic operations.
// Controller defines the business logic methods for authentication.
type Controller interface {
	Register(ctx context.Context, input RegisterInput) (model.User, error)
	Login(ctx context.Context, input LoginInput) (model.User, string, error)
	GetProfile(ctx context.Context, userID int64) (model.User, error)
}

type impl struct {
	repo repository.Registry
}

// New creates a new authentication controller instance.
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
