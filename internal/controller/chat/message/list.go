package message

import (
	"context"
)

// List retrieves a paginated list of messages for a channel.
func (i impl) List(ctx context.Context, channelID int64, limit int, offset int) ([]MessageResponse, error) {
	messages, err := i.repo.Message().ListByChannel(ctx, channelID, limit, offset)
	if err != nil {
		return nil, err
	}

	res := make([]MessageResponse, len(messages))
	for idx, m := range messages {
		res[idx] = MessageResponse{
			ID:        m.ID,
			ChannelID: m.ChannelID,
			SenderID:  m.SenderID,
			ParentID:  m.ParentID,
			Content:   m.Content,
			IsEdited:  m.IsEdited,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return res, nil
}
