package organizations

import (
	"context"
	"errors"

	"asona/internal/repository/organizations"

	pkgerrors "github.com/pkg/errors"
)

// Get retrieves an organization's details by its ID.
func (i impl) Get(ctx context.Context, oid int64) (OrganizationResponse, error) {
	org, err := i.repo.Organization().GetByID(ctx, oid)
	if err != nil {
		if errors.Is(err, organizations.ErrOrganizationNotFound) {
			return OrganizationResponse{}, pkgerrors.WithStack(ErrOrganizationNotFound)
		}
		return OrganizationResponse{}, err
	}

	return OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		LogoURL:     org.LogoURL,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}
