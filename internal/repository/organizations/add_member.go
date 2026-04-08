package organizations

import (
	"asona/internal/model"
	"context"
)

func (i impl) AddMember(ctx context.Context, member model.OrganizationMember) (model.OrganizationMember, error) {
	return member, nil
}
