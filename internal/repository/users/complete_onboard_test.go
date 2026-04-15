package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeContextExecutor struct {
	execFn func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (f fakeContextExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if f.execFn != nil {
		return f.execFn(ctx, query, args...)
	}
	return nil, nil
}

func (f fakeContextExecutor) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f fakeContextExecutor) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return nil
}

type fakeSQLResult struct {
	rowsAffected int64
	rowsErr      error
}

func (f fakeSQLResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (f fakeSQLResult) RowsAffected() (int64, error) {
	if f.rowsErr != nil {
		return 0, f.rowsErr
	}
	return f.rowsAffected, nil
}

func TestCompleteOnboard(t *testing.T) {
	ctx := context.Background()
	userID := int64(12)

	t.Run("success when one row updated", func(t *testing.T) {
		repo := impl{db: fakeContextExecutor{execFn: func(context.Context, string, ...interface{}) (sql.Result, error) {
			return fakeSQLResult{rowsAffected: 1}, nil
		}}}

		err := repo.CompleteOnboard(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("returns wrapped exec error", func(t *testing.T) {
		execErr := errors.New("exec failed")
		repo := impl{db: fakeContextExecutor{execFn: func(context.Context, string, ...interface{}) (sql.Result, error) {
			return nil, execErr
		}}}

		err := repo.CompleteOnboard(ctx, userID)
		require.Error(t, err)
		require.ErrorIs(t, err, execErr)
		require.Contains(t, err.Error(), "failed to complete onboard")
	})

	t.Run("returns wrapped rows-affected error", func(t *testing.T) {
		rowsErr := errors.New("rows affected failed")
		repo := impl{db: fakeContextExecutor{execFn: func(context.Context, string, ...interface{}) (sql.Result, error) {
			return fakeSQLResult{rowsErr: rowsErr}, nil
		}}}

		err := repo.CompleteOnboard(ctx, userID)
		require.Error(t, err)
		require.ErrorIs(t, err, rowsErr)
		require.Contains(t, err.Error(), "failed to get rows affected")
	})

	t.Run("returns user-not-found when no rows updated", func(t *testing.T) {
		repo := impl{db: fakeContextExecutor{execFn: func(context.Context, string, ...interface{}) (sql.Result, error) {
			return fakeSQLResult{rowsAffected: 0}, nil
		}}}

		err := repo.CompleteOnboard(ctx, userID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserNotFound)
	})
}
