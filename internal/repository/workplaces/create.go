package workplaces

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"asona/internal/model"
)

// Create inserts a new workplace record into the database.
func (i impl) Create(ctx context.Context, wp model.Workplace) (model.Workplace, error) {
	query := `
		INSERT INTO public.workplaces (name, icon_url, size, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := i.db.QueryRowContext(ctx, query, wp.Name, wp.IconURL, wp.Size, wp.CreatedBy).
		Scan(&wp.ID, &wp.CreatedAt, &wp.UpdatedAt)
	if err != nil {
		return model.Workplace{}, pkgerrors.WithStack(fmt.Errorf("failed to insert workplace: %w", err))
	}

	return wp, nil
}
