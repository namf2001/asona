package channel

import (
	"context"
	"time"

	"asona/internal/model"
)

// CreateChannelInput represents the input data for creating a new channel.
type CreateChannelInput struct {
	WorkplaceID int64  `json:"workplace_id"`
	ProjectID   int64  `json:"project_id,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

// ChannelResponse represents the public channel data.
type ChannelResponse struct {
	ID          int64     `json:"id"`
	WorkplaceID int64     `json:"workplace_id"`
	ProjectID   int64     `json:"project_id,omitempty"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create creates a new chat channel and adds the creator as a member.
func (i impl) Create(ctx context.Context, userID int64, input CreateChannelInput) (ChannelResponse, error) {
	ch := model.Channel{
		WorkplaceID: input.WorkplaceID,
		ProjectID:   input.ProjectID,
		Name:        input.Name,
		Type:        model.ChannelType(input.Type),
		CreatedBy:   userID,
	}

	// Gọi trực tiếp từ Registry
	created, err := i.repo.Channel().Create(ctx, ch)
	if err != nil {
		return ChannelResponse{}, err
	}

	err = i.repo.Channel().AddMember(ctx, created.ID, userID)
	if err != nil {
		return ChannelResponse{}, err
	}

	return ChannelResponse{
		ID:          created.ID,
		WorkplaceID: created.WorkplaceID,
		ProjectID:   created.ProjectID,
		Name:        created.Name,
		Type:        string(created.Type),
		CreatedBy:   created.CreatedBy,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}, nil
}
