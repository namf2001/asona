package organizations

import (
	"asona/internal/model"
	"context"
)

func (i impl) GetByID(ctx context.Context, id int64) (model.Organization, error) {
	return model.Organization{}, nil
}
