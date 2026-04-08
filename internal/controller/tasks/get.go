package tasks

import (
	"context"
)

// GetByID retrieves a single task's details by its unique identifier.
func (i impl) GetByID(ctx context.Context, id int64) (TaskResponse, error) {
	t, err := i.repo.Task().GetByID(ctx, id)
	if err != nil {
		return TaskResponse{}, err
	}

	return TaskResponse{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		CreatedBy:   t.CreatedBy,
		AssigneeID:  t.AssigneeID,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}, nil
}
