package users

import (
	"asona/internal/model"
	"context"
	"fmt"
	pkgerrors "github.com/pkg/errors"
)

func (i impl) Create(ctx context.Context, user model.User) (model.User, error) {
	user.Prepare()
	query := `
		INSERT INTO public.users (name, username, display_name, email, password, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err := i.db.QueryRowContext(ctx, query,
		user.Name,
		user.Username,
		user.DisplayName,
		user.Email,
		user.Password,
		user.Image,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return model.User{}, pkgerrors.WithStack(fmt.Errorf("failed to create user: %w", err))
	}

	return user, nil
}
