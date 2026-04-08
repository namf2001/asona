package sessions

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// Delete removes a session by token.
func (i impl) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM public.sessions WHERE session_token = $1`

	result, err := i.db.ExecContext(ctx, query, token)
	if err != nil {
		return pkgerrors.WithStack(fmt.Errorf("failed to delete session: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkgerrors.WithStack(fmt.Errorf("failed to read delete result: %w", err))
	}
	if rowsAffected == 0 {
		return pkgerrors.WithStack(ErrSessionNotFound)
	}

	return nil
}
