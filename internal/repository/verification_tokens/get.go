package verification_tokens

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"asona/internal/model"

	pkgerrors "github.com/pkg/errors"
)

func (i impl) GetValidToken(ctx context.Context, identifier, tokenStr string) (model.VerificationToken, error) {
	const query = `
		SELECT identifier, token, expires
		FROM verification_token
		WHERE identifier = $1 AND token = $2 AND expires > $3
	`
	var token model.VerificationToken
	err := i.db.QueryRowContext(ctx, query, identifier, tokenStr, time.Now()).Scan(
		&token.Identifier,
		&token.Token,
		&token.Expires,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.VerificationToken{}, pkgerrors.WithStack(ErrTokenNotFoundOrExpired)
		}
		return model.VerificationToken{}, pkgerrors.WithStack(err)
	}
	return token, nil
}
