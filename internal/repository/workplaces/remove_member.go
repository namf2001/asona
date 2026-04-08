package workplaces

import (
	"context"
)

// RemoveMember deletes a user's membership from a workplace.
func (i impl) RemoveMember(ctx context.Context, workplaceID, userID int64) error {
	return nil
}
