package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/cenkalti/backoff/v4"
	"asona/internal/pkg/logger"
	pkgerrors "github.com/pkg/errors"
)

// Tx starts a transaction with default backoff policy (3 retries, 1 minute max)
func Tx(ctx context.Context, dbconn BeginnerExecutor, callback func(ContextExecutor) error) error {
	return TxWithBackOff(ctx, ExponentialBackOff(3, time.Minute), dbconn, callback)
}

// TxWithBackOff starts a transaction with the provided backoff policy.
// It handles begin, commit, and rollback automatically.
// If callback returns an error, the transaction is rolled back.
// If callback succeeds, the transaction is committed.
func TxWithBackOff(ctx context.Context, b backoff.BackOff, dbconn BeginnerExecutor, callback func(ContextExecutor) error) error {
	if b == nil {
		b = &backoff.StopBackOff{}
	}

	tx, err := beginTx(ctx, dbconn, b)
	if err != nil {
		return err
	}

	var committed bool
	defer func() {
		if committed {
			return
		}
		// Always rollback if not committed to clean up
		_ = tx.Rollback()
	}()

	// Execute the callback within the transaction
	if err = callback(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return pkgerrors.WithStack(err)
	}

	committed = true
	return nil
}

func beginTx(ctx context.Context, dbconn BeginnerExecutor, b backoff.BackOff) (*sql.Tx, error) {
	var tryCount int
	var tx *sql.Tx
	if err := backoff.Retry(func() error {
		tryCount++
		var err error

		logger.INFO.Printf("DB: BeginTx Attempt: %d", tryCount)
		tx, err = dbconn.BeginTx(ctx, nil)

		return pkgerrors.WithStack(err)
	}, backoff.WithContext(b, ctx)); err != nil {
		return nil, err
	}
	return tx, nil
}
