package auth

import (
	"context"

	"asona/internal/model"
	"asona/internal/repository"
	"asona/internal/service/oauth"
)

// Controller defines the business logic methods for authentication.
type Controller interface {
	Register(ctx context.Context, input RegisterInput) (model.User, error)
	Login(ctx context.Context, input LoginInput) (model.User, string, error)
	GetProfile(ctx context.Context, userID int64) (model.User, error)
	Logout(ctx context.Context, token string) error
	GoogleAuthURL(ctx context.Context, state string) (string, error)
	GoogleCallback(ctx context.Context, code string) (model.User, string, error)
}

type impl struct {
	repo  repository.Registry
	oauth oauth.Service
}

// New creates a new authentication controller instance.
func New(repo repository.Registry, oauthSvc oauth.Service) Controller {
	return impl{repo: repo, oauth: oauthSvc}
}
