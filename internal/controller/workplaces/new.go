package workplaces

import (
	"context"

	"asona/internal/repository"
)

// Controller defines the business logic methods for workplace management.
type Controller interface {
	Create(ctx context.Context, userID int64, input CreateWorkplaceInput) (WorkplaceResponse, error)
}

type impl struct {
	repo repository.Registry
}

// New creates a new workplaces controller instance.
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
