package verification_tokens

import (
	"context"

	pkgerrors "github.com/pkg/errors"
)

func (i impl) Delete(ctx context.Context, identifier, tokenStr string) error {
	const query = `DELETE FROM verification_token WHERE identifier = $1 AND token = $2`
	_, err := i.db.ExecContext(ctx, query, identifier, tokenStr)
	return pkgerrors.WithStack(err)
}

func (i impl) DeleteAllForIdentifier(ctx context.Context, identifier string) error {
	const query = `DELETE FROM verification_token WHERE identifier = $1`
	_, err := i.db.ExecContext(ctx, query, identifier)
	if err != nil {
		return pkgerrors.WithStack(err)
	}
	return nil
}
