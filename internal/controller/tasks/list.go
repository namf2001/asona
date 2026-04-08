package tasks

import (
	"context"
	"asona/internal/repository/tasks"
)

// ListByProject retrieves all tasks belonging to a specific project.
func (i impl) ListByProject(ctx context.Context, projectID int64) ([]TaskResponse, error) {
	list, err := i.repo.Task().List(ctx, tasks.ListFilter{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}

	res := make([]TaskResponse, 0, len(list))
	for _, t := range list {
		res = append(res, TaskResponse{
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
		})
	}

	return res, nil
}
