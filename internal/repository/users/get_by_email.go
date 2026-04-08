package users

import (
	"asona/internal/model"
	"context"
	"fmt"
)

func (i impl) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `
		SELECT id, name, username, display_name, email, password, avatar_url, is_active, created_at, updated_at
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
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}
