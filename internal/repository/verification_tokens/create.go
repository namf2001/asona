package verification_tokens

import (
	"context"

	"asona/internal/model"

	pkgerrors "github.com/pkg/errors"
)

func (i impl) Create(ctx context.Context, token model.VerificationToken) error {
	const query = `
		INSERT INTO verification_token (identifier, token, expires)
		VALUES ($1, $2, $3)
	`
	_, err := i.db.ExecContext(ctx, query, token.Identifier, token.Token, token.Expires)
	if err != nil {
		return pkgerrors.WithStack(err)
	}
	return nil
}
