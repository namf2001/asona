package channels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"asona/internal/model"
)

// GetByID retrieves a specific channel's details by its ID.
func (i impl) GetByID(ctx context.Context, id int64) (model.Channel, error) {
	var c model.Channel
	query := `
		SELECT id, workplace_id, project_id, name, type, created_by, created_at, updated_at
		FROM public.channels
		WHERE id = $1
	`
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.WorkplaceID,
		&c.ProjectID,
		&c.Name,
		&c.Type,
		&c.CreatedBy,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Channel{}, ErrChannelNotFound
		}
		return model.Channel{}, fmt.Errorf("failed to get channel: %w", err)
	}

	return c, nil
}
