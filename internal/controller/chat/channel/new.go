package channel

import (
	"context"

	"asona/internal/repository"
)

// Controller defines channel business logic operations.
type Controller interface {
	// Create creates a new channel and adds the creator as first member
	Create(ctx context.Context, userID int64, input CreateChannelInput) (ChannelResponse, error)
	// GetByID retrieves channel details by ID
	GetByID(ctx context.Context, channelID int64) (ChannelResponse, error)
	// AddMember adds a user to a channel
	AddMember(ctx context.Context, channelID, userID int64) error
}

type impl struct {
	repo repository.Registry
}

// New returns a new channel Controller instance.
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
