package pg

import (
	"context"
	"database/sql"
)

// ContextExecutor can perform SQL queries with context.
// Both *sql.DB and *sql.Tx satisfy this interface.
type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Beginner allows creation of context aware transactions with options.
type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// BeginnerExecutor can context-aware perform SQL queries and
// create context-aware transactions with options.
type BeginnerExecutor interface {
	Beginner
	ContextExecutor

	PingContext(ctx context.Context) error
	Close() error
}
