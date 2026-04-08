package users

import (
	"asona/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

func (i impl) GetByID(ctx context.Context, id int64) (model.User, error) {
	var user model.User
	query := `
		SELECT id, name, username, display_name, email, password, avatar_url, is_active, created_at, updated_at
		FROM public.users
		WHERE id = $1 AND is_active = true
	`
	err := i.db.QueryRowContext(ctx, query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, pkgerrors.WithStack(ErrUserNotFound)
		}
		return model.User{}, pkgerrors.WithStack(fmt.Errorf("failed to get user by id: %w", err))
	}

	return user, nil
}
