package auth

import "asona/internal/repository"

// Controller defines the business logic for authenticated operations.
type Controller interface {
}

type impl struct {
	repo repository.Registry
}

// New returns a new authenticated Controller.
func New(repo repository.Registry) Controller {
	return impl{
		repo: repo,
	}
}
