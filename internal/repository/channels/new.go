package channels

import (
	"asona/internal/model"
	"asona/internal/repository/db/pg"
	"context"
)

// Repository defines the interface for channel data operations.
type Repository interface {
	// Create inserts a new channel into the database.
	Create(ctx context.Context, input model.Channel) (model.Channel, error)
	// GetByID retrieves a channel by its unique identifier.
	GetByID(ctx context.Context, id int64) (model.Channel, error)
	// AddMember adds a user to a specific channel.
	AddMember(ctx context.Context, channelID int64, userID int64) error
}

type impl struct {
	db pg.ContextExecutor
}

// New creates a new instance of channel repository.
func New(db pg.ContextExecutor) Repository {
	return impl{
		db: db,
	}
}
