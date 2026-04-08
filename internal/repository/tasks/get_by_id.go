package tasks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"asona/internal/model"
)

// GetByID retrieves a single task by its unique identifier.
func (i impl) GetByID(ctx context.Context, id int64) (model.Task, error) {
	var t model.Task
	query := `
		SELECT id, project_id, title, description, status, priority, created_by, assignee_id, due_date, created_at, updated_at
		FROM public.tasks
		WHERE id = $1
	`
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.ProjectID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.CreatedBy,
		&t.AssigneeID,
		&t.DueDate,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, ErrTaskNotFound
		}
		return model.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return t, nil
}
