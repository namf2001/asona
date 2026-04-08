package organizations

import (
	"context"
	"time"

	"asona/internal/model"
	"asona/internal/repository"
)

// CreateOrganizationInput represents the input data for creating an organization.
type CreateOrganizationInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

// OrganizationResponse represents the public organization data.
type OrganizationResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create generates a new organization in the system and adds the creator as the initial owner.
func (i impl) Create(ctx context.Context, userID int64, input CreateOrganizationInput) (OrganizationResponse, error) {
	var res OrganizationResponse

	err := i.repo.DoInTx(ctx, func(ctx context.Context, txRepo repository.Registry) error {
		// 1. Create organization
		org, err := txRepo.Organization().Create(ctx, model.Organization{
			Name:        input.Name,
			Description: input.Description,
			LogoURL:     input.LogoURL,
			CreatedBy:   userID,
		})
		if err != nil {
			return err
		}

		// 2. Add creator as admin
		_, err = txRepo.Organization().AddMember(ctx, model.OrganizationMember{
			OrganizationID: org.ID,
			UserID:         userID,
			Role:           model.OrgRoleAdmin,
		})
		if err != nil {
			return err
		}

		res = OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			LogoURL:     org.LogoURL,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		}
		return nil
	}, nil)

	return res, err
}
