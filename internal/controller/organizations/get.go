package organizations

import (
	"context"
)

// Get retrieves an organization's details by its ID.
func (i impl) Get(ctx context.Context, oid int64) (OrganizationResponse, error) {
	return OrganizationResponse{}, nil
}
