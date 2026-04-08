package organizations

import (
	"context"
)

// Get retrieves an organization's details by its ID.
func (i impl) Get(ctx context.Context, oid int64) (OrganizationResponse, error) {
	org, err := i.repo.Organization().GetByID(ctx, oid)
	if err != nil {
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
