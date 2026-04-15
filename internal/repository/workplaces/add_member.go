package workplaces

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

// AddMember registers a user as a member of a workplace.
func (i impl) AddMember(ctx context.Context, member model.WorkplaceMember) error {
	query := `
		INSERT INTO public.workplace_members (workplace_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (workplace_id, user_id) DO NOTHING
	`
	_, err := i.db.ExecContext(ctx, query, member.WorkplaceID, member.UserID, member.Role)
	if err != nil {
		return pkgerrors.WithStack(fmt.Errorf("failed to add workplace member: %w", err))
	}

	return nil
}
