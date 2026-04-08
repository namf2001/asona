package channel

import (
	"context"
)

// AddMember adds a existing user to a chat channel.
func (i impl) AddMember(ctx context.Context, channelID int64, userID int64) error {
	return i.repo.Channel().AddMember(ctx, channelID, userID)
}
