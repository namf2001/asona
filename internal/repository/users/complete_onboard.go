package users

import (
	"context"
	"fmt"
	"time"

	pkgerrors "github.com/pkg/errors"
)

// CompleteOnboard sets the onboarded_at timestamp for the given user,
// marking them as having completed the onboarding flow.
func (i impl) CompleteOnboard(ctx context.Context, userID int64) error {
	query := `
		UPDATE public.users
		SET onboarded_at = $1,
			onboarding_status = 'completed',
			onboarding_step = 3,
			updated_at = $2
		WHERE id = $3
	`
	now := time.Now()
	result, err := i.db.ExecContext(ctx, query, now, now, userID)
	if err != nil {
		return pkgerrors.WithStack(fmt.Errorf("failed to complete onboard for user %d: %w", userID, err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkgerrors.WithStack(fmt.Errorf("failed to get rows affected while completing onboard for user %d: %w", userID, err))
	}

	if rowsAffected == 0 {
		return pkgerrors.WithStack(ErrUserNotFound)
	}

	return nil
}
