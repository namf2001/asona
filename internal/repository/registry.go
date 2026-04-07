package repository

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	pkgerrors "github.com/pkg/errors"

	"asona/internal/repository/db/pg"
)

// Registry is the registry of all domain-specific repositories
// and provides transaction capabilities.
type Registry interface {
	// DoInTx wraps operations within a Postgres transaction.
	DoInTx(ctx context.Context, txFunc func(ctx context.Context, txRepo Registry) error, overrideBackoffPolicy backoff.BackOff) error
}

// New returns a new instance of Registry backed by a Postgres connection.
func New(db pg.BeginnerExecutor) Registry {
	return impl{
		pgConn: db,
	}
}

type impl struct {
	pgConn pg.BeginnerExecutor // Used to start DB transactions
	tx     pg.ContextExecutor  // Non-nil when inside a transaction
}

// DoInTx wraps operations within a Postgres transaction with backoff retry.
// Nested transactions are not allowed.
func (i impl) DoInTx(
	ctx context.Context,
	txFunc func(ctx context.Context, txRepo Registry) error,
	overrideBackoffPolicy backoff.BackOff,
) error {
	if i.tx != nil {
		return pkgerrors.WithStack(errNestedTx)
	}

	if overrideBackoffPolicy == nil {
		overrideBackoffPolicy = pg.ExponentialBackOff(3, time.Minute)
	}

	return pg.TxWithBackOff(ctx, overrideBackoffPolicy, i.pgConn, func(tx pg.ContextExecutor) error {
		newI := impl{
			tx: tx,
		}
		return txFunc(ctx, newI)
	})
}
