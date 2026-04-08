package channels

import (
	"asona/internal/model"
	"context"
)

// GetByID retrieves a channel by its unique identifier.
func (i impl) GetByID(ctx context.Context, id int64) (model.Channel, error) {
	return model.Channel{}, nil
}
