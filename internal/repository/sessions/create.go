package sessions

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

// Create inserts a new session row for an authenticated user.
func (i impl) Create(ctx context.Context, session model.Session) (model.Session, error) {
	query := `
		INSERT INTO public.sessions (user_id, session_token, expires_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::inet)
		RETURNING id, created_at
	`

	err := i.db.QueryRowContext(ctx, query,
		session.UserID,
		session.SessionToken,
		session.ExpiresAt,
		session.UserAgent,
		session.IPAddress,
	).Scan(&session.ID, &session.CreatedAt)
	if err != nil {
		return model.Session{}, pkgerrors.WithStack(fmt.Errorf("failed to create session: %w", err))
	}

	return session, nil
}
