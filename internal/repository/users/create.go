package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			// Distinguish which unique constraint was violated by its name.
			switch pqErr.Constraint {
			case "users_username_key":
				return model.User{}, pkgerrors.WithStack(ErrUsernameAlreadyExists)
			default:
				// email or any other unique column
				return model.User{}, pkgerrors.WithStack(ErrEmailAlreadyExists)
			}
		}
		return model.User{}, pkgerrors.WithStack(fmt.Errorf("failed to create user: %w", err))
	}

	return user, nil
}
