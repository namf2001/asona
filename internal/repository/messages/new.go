package messages

import (
	"context"

	"asona/internal/model"
	"asona/internal/repository/db/pg"
)

// Repository defines message data access operations
type Repository interface {
	// Create persists a new message and returns the created record
	Create(ctx context.Context, message model.Message) (model.Message, error)
	// ListByChannel retrieves paginated messages for a channel
	ListByChannel(ctx context.Context, channelID int64, limit, offset int) ([]model.Message, error)
}

type impl struct {
	db pg.ContextExecutor
}

// New returns a new message Repository instance.
func New(db pg.ContextExecutor) Repository {
	return impl{
		db: db,
	}
}
