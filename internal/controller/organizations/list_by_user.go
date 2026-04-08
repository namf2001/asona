package organizations

import (
	"context"
	"time"
	"asona/internal/model"
)

// OrganizationWithRoleResponse represents an organization including the user's role.
type OrganizationWithRoleResponse struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	LogoURL     string        `json:"logo_url"`
	Role        model.OrgRole `json:"role"`
	JoinedAt    time.Time     `json:"joined_at"`
}

// ListByUser returns all organizations a specific user belongs to.
func (i impl) ListByUser(ctx context.Context, userID int64) ([]OrganizationWithRoleResponse, error) {
	return nil, nil
}
