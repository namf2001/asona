package channels

import (
	"asona/internal/model"
	"context"
)

// Create inserts a new channel into the database.
func (i impl) Create(ctx context.Context, input model.Channel) (model.Channel, error) {
	return model.Channel{}, nil
}
