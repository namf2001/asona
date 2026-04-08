package organizations

import (
	"asona/internal/repository"
	"context"
)

// Controller defines the business logic methods for organization management.
type Controller interface {
	Create(ctx context.Context, userID int64, input CreateOrganizationInput) (OrganizationResponse, error)
	Get(ctx context.Context, orgID int64) (OrganizationResponse, error)
	ListByUser(ctx context.Context, userID int64) ([]OrganizationWithRoleResponse, error)
}

type impl struct {
	repo repository.Registry
}

// New creates a new organizations controller instance.
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
