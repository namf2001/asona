package sessions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

// GetByToken retrieves an active (not expired) session by session token.
func (i impl) GetByToken(ctx context.Context, token string) (model.Session, error) {
	var session model.Session
	var userAgent sql.NullString
	var ipAddress sql.NullString

	query := `
		SELECT id, user_id, session_token, expires_at, user_agent, ip_address::text, created_at
		FROM public.sessions
		WHERE session_token = $1 AND expires_at > NOW()
	`

	err := i.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionToken,
		&session.ExpiresAt,
		&userAgent,
		&ipAddress,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Session{}, pkgerrors.WithStack(ErrSessionNotFound)
		}
		return model.Session{}, pkgerrors.WithStack(fmt.Errorf("failed to get session by token: %w", err))
	}

	if userAgent.Valid {
		session.UserAgent = userAgent.String
	}
	if ipAddress.Valid {
		session.IPAddress = ipAddress.String
	}

	return session, nil
}
