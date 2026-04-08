package messages

import (
	"asona/internal/model"
	"context"
)

// ListByChannel retrieves paginated messages for a channel.
func (i impl) ListByChannel(ctx context.Context, channelID int64, limit int, offset int) ([]model.Message, error) {
	return nil, nil
}
