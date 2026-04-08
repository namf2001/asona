package tasks

import (
	"context"
	"time"

	"asona/internal/model"
)

// UpdateTaskInput represents the data for updating an existing task.
type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssigneeID  int64  `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

// Update updates an existing task's information.
func (i impl) Update(ctx context.Context, senderID int64, taskID int64, input UpdateTaskInput) error {
	return i.repo.Task().Update(ctx, model.Task{
		ID:          taskID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    model.TaskPriority(input.Priority),
		AssigneeID:  input.AssigneeID,
		DueDate:     input.DueDate,
	})
}
