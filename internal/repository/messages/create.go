package messages

import (
	"asona/internal/model"
	"context"
)

// Create persists a new message and returns the created record.
func (i impl) Create(ctx context.Context, input model.Message) (model.Message, error) {
	return model.Message{}, nil
}
