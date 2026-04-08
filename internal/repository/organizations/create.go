package organizations

import (
	"context"
	"fmt"
	"time"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

func (i impl) Create(ctx context.Context, org model.Organization) (model.Organization, error) {
	now := time.Now()
	org.CreatedAt = now
	org.UpdatedAt = now

	query := `
		INSERT INTO public.organizations (name, description, logo_url, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := i.db.QueryRowContext(ctx, query,
		org.Name,
		org.Description,
		org.LogoURL,
		org.CreatedBy,
		org.CreatedAt,
		org.UpdatedAt,
	).Scan(&org.ID)

	if err != nil {
		return model.Organization{}, pkgerrors.WithStack(fmt.Errorf("failed to create organization: %w", err))
	}

	return org, nil
}
