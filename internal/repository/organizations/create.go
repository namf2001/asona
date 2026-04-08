package organizations

import (
	"asona/internal/model"
	"context"
)

func (i impl) Create(ctx context.Context, org model.Organization) (model.Organization, error) {
	return org, nil
}
