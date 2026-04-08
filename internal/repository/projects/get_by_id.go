package projects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"asona/internal/model"
)

// GetByID retrieves a specific project by its ID.
func (i impl) GetByID(ctx context.Context, id int64) (model.Project, error) {
	var p model.Project
	query := `
		SELECT id, workplace_id, name, description, created_by, created_at, updated_at
		FROM public.projects
		WHERE id = $1
	`
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.WorkplaceID,
		&p.Name,
		&p.Description,
		&p.CreatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Project{}, ErrProjectNotFound
		}
		return model.Project{}, fmt.Errorf("failed to get project: %w", err)
	}

	return p, nil
}
