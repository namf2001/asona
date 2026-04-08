package organizations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

func (i impl) GetByID(ctx context.Context, id int64) (model.Organization, error) {
	var org model.Organization
	query := `
		SELECT id, name, description, logo_url, created_by, created_at, updated_at
		FROM public.organizations
		WHERE id = $1
	`
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Description,
		&org.LogoURL,
		&org.CreatedBy,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Organization{}, pkgerrors.WithStack(ErrOrganizationNotFound)
		}
		return model.Organization{}, pkgerrors.WithStack(fmt.Errorf("failed to get organization: %w", err))
	}

	return org, nil
}
