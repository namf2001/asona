package users

import (
	"asona/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

func (i impl) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `
		SELECT id, name, username, display_name, email, password, avatar_url, is_active, onboarding_status, onboarding_step, onboarded_at, created_at, updated_at
		FROM public.users
		WHERE email = $1 AND is_active = true
	`
	err := i.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.IsActive,
		&user.OnboardingStatus,
		&user.OnboardingStep,
		&user.OnboardedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, pkgerrors.WithStack(ErrUserNotFound)
		}
		return model.User{}, pkgerrors.WithStack(fmt.Errorf("failed to get user by email: %w", err))
	}

	return user, nil
}
