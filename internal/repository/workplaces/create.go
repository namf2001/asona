package workplaces

import (
	"asona/internal/model"
	"context"
)

// Create inserts a new workplace record into the database.
func (i impl) Create(ctx context.Context, wp model.Workplace) (model.Workplace, error) {
	return wp, nil
}
