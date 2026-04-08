package message

import (
	"context"

	"asona/internal/repository"
)

// Controller defines message business logic operations.
type Controller interface {
	// Send creates a new message in a channel
	Send(ctx context.Context, senderID int64, input SendMessageInput) (MessageResponse, error)
	// List retrieves paginated messages for a channel
	List(ctx context.Context, channelID int64, limit, offset int) ([]MessageResponse, error)
}

type impl struct {
	repo repository.Registry
}

// New returns a new message Controller instance.
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
