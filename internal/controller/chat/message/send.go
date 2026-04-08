package message

import (
	"context"
	"time"

	"asona/internal/model"
)

// SendMessageInput represents the input data for sending a new message.
type SendMessageInput struct {
	ChannelID int64  `json:"channel_id"`
	ParentID  int64  `json:"parent_id,omitempty"`
	Content   string `json:"content"`
}

// MessageResponse represents the public message data.
type MessageResponse struct {
	ID        int64     `json:"id"`
	ChannelID int64     `json:"channel_id"`
	SenderID  int64     `json:"sender_id"`
	ParentID  int64     `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"is_edited"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Send processes and saves a new message sent by a user.
func (i impl) Send(ctx context.Context, userID int64, input SendMessageInput) (MessageResponse, error) {
	msg := model.Message{
		ChannelID: input.ChannelID,
		SenderID:  userID,
		ParentID:  input.ParentID,
		Content:   input.Content,
	}
	// Gọi trực tiếp từ Registry
	created, err := i.repo.Message().Create(ctx, msg)
	if err != nil {
		return MessageResponse{}, err
	}

	return MessageResponse{
		ID:        created.ID,
		ChannelID: created.ChannelID,
		SenderID:  created.SenderID,
		ParentID:  created.ParentID,
		Content:   created.Content,
		IsEdited:  created.IsEdited,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}
