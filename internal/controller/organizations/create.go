package organizations

import (
	"context"
	"time"
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

// Create generates a new organization in the system.
func (i impl) Create(ctx context.Context, userID int64, input CreateOrganizationInput) (OrganizationResponse, error) {
	return OrganizationResponse{}, nil
}
