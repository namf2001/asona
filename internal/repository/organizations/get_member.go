package organizations

import (
	"asona/internal/model"
	"context"
)

func (i impl) GetMember(ctx context.Context, orgID, userID int64) (model.OrganizationMember, error) {
	return model.OrganizationMember{}, nil
}
