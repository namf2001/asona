package workplaces

import (
	"asona/internal/model"
	"context"
)

// GetByID retrieves a single workplace by its unique identifier.
func (i impl) GetByID(ctx context.Context, id int64) (model.Workplace, error) {
	return model.Workplace{}, nil
}
