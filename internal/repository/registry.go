package repository

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	pkgerrors "github.com/pkg/errors"

	"asona/internal/repository/accounts"
	repoChannel "asona/internal/repository/channels"
	"asona/internal/repository/db/pg"
	repoMessage "asona/internal/repository/messages"
	"asona/internal/repository/organizations"
	"asona/internal/repository/projects"
	"asona/internal/repository/sessions"
	"asona/internal/repository/tasks"
	"asona/internal/repository/users"
	"asona/internal/repository/workplaces"
)

// Registry is the registry of all domain-specific repositories
// and provides transaction capabilities.
type Registry interface {
	// User return user repository
	User() users.Repository
	// Account return account repository
	Account() accounts.Repository
	// Session return session repository
	Session() sessions.Repository
	// Organization return organization repository
	Organization() organizations.Repository
	// Workplace return workplace repository
	Workplace() workplaces.Repository
	// Project return project repository
	Project() projects.Repository
	// Task return task repository
	Task() tasks.Repository
	// Channel return channel repository
	Channel() repoChannel.Repository
	// Message return message repository
	Message() repoMessage.Repository
	// DoInTx wraps operations within a Postgres transaction.
	DoInTx(ctx context.Context, txFunc func(ctx context.Context, txRepo Registry) error, overrideBackoffPolicy backoff.BackOff) error
}

// New returns a new instance of Registry backed by a Postgres connection.
func New(db pg.BeginnerExecutor) Registry {
	return impl{
		pgConn:        db,
		users:         users.New(db),
		accounts:      accounts.New(db),
		sessions:      sessions.New(db),
		organizations: organizations.New(db),
		workplaces:    workplaces.New(db),
		projects:      projects.New(db),
		tasks:         tasks.New(db),
		channels:      repoChannel.New(db),
		messages:      repoMessage.New(db),
	}
}

type impl struct {
	pgConn        pg.BeginnerExecutor // Used to start DB transactions
	tx            pg.ContextExecutor  // Non-nil when inside a transaction
	users         users.Repository
	accounts      accounts.Repository
	sessions      sessions.Repository
	organizations organizations.Repository
	workplaces    workplaces.Repository
	projects      projects.Repository
	tasks         tasks.Repository
	channels      repoChannel.Repository
	messages      repoMessage.Repository
}

func (i impl) User() users.Repository {
	return i.users
}

func (i impl) Account() accounts.Repository {
	return i.accounts
}

func (i impl) Session() sessions.Repository {
	return i.sessions
}

func (i impl) Organization() organizations.Repository {
	return i.organizations
}

func (i impl) Workplace() workplaces.Repository {
	return i.workplaces
}

func (i impl) Project() projects.Repository {
	return i.projects
}

func (i impl) Task() tasks.Repository {
	return i.tasks
}

func (i impl) Channel() repoChannel.Repository {
	return i.channels
}

func (i impl) Message() repoMessage.Repository {
	return i.messages
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
			tx:            tx,
			users:         users.New(tx),
			accounts:      accounts.New(tx),
			sessions:      sessions.New(tx),
			organizations: organizations.New(tx),
			workplaces:    workplaces.New(tx),
			projects:      projects.New(tx),
			tasks:         tasks.New(tx),
			channels:      repoChannel.New(tx),
			messages:      repoMessage.New(tx),
		}
		return txFunc(ctx, newI)
	})
}
