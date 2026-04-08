package organizations

import (
	"asona/internal/model"
	"asona/internal/repository/db/pg"
	"context"
)

type OrganizationWithRole struct {
	Organization model.Organization `json:"organization"`
	Role         model.OrgRole      `json:"role"`
	JoinedAt     string             `json:"joined_at"`
}

type Repository interface {
	// Create creates a new organization
	Create(ctx context.Context, org model.Organization) (model.Organization, error)

	// GetByID retrieves an organization by ID
	GetByID(ctx context.Context, id int64) (model.Organization, error)

	// AddMember adds a member to organization
	AddMember(ctx context.Context, member model.OrganizationMember) (model.OrganizationMember, error)

	// GetMember gets a specific member
	GetMember(ctx context.Context, orgID, userID int64) (model.OrganizationMember, error)

	// RemoveMember removes a member from organization
	RemoveMember(ctx context.Context, orgID, userID int64) error
}

type impl struct {
	db pg.ContextExecutor
}

func New(db pg.ContextExecutor) Repository {
	return impl{
		db: db,
	}
}
