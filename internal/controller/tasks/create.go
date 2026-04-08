package tasks

import (
	"context"
	"time"

	"asona/internal/model"
)

// CreateTaskInput represents the data required to create a new task.
type CreateTaskInput struct {
	ProjectID   int64  `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssigneeID  int64  `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

// TaskResponse represents the public task data.
type TaskResponse struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedBy   int64     `json:"created_by"`
	AssigneeID  int64     `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create creates a new task in the specified project.
func (i impl) Create(ctx context.Context, senderID int64, input CreateTaskInput) (TaskResponse, error) {
	t, err := i.repo.Task().Create(ctx, model.Task{
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    model.TaskPriority(input.Priority),
		CreatedBy:   senderID,
		AssigneeID:  input.AssigneeID,
		DueDate:     input.DueDate,
	})
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
