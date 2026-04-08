package channel

import (
	"context"
)

// GetByID retrieves a specific channel's details by its ID.
func (i impl) GetByID(ctx context.Context, id int64) (ChannelResponse, error) {
	created, err := i.repo.Channel().GetByID(ctx, id)
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
