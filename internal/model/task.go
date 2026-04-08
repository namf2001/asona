package model

import (
	"encoding/json"
	"time"
)

// Task represents a work item
type Task struct {
	ID          int64        `json:"id,omitempty"          db:"id"`
	ProjectID   int64        `json:"project_id,omitempty"  db:"project_id"`
	Status      string       `json:"status,omitempty"      db:"status"`
	AssigneeID  int64        `json:"assignee_id,omitempty" db:"assignee_id"`
	ParentID    int64        `json:"parent_id,omitempty"   db:"parent_id"` // Subtask support
	Title       string       `json:"title,omitempty"       db:"title"`
	Description string       `json:"description,omitempty" db:"description"`
	Priority    TaskPriority `json:"priority,omitempty"    db:"priority"`
	Position    float64      `json:"position,omitempty"    db:"position"`
	DueDate     *time.Time   `json:"due_date,omitempty"    db:"due_date"`
	StartDate   *time.Time   `json:"start_date,omitempty"  db:"start_date"`
	Estimate    int          `json:"estimate,omitempty"    db:"estimate"` // In minutes/hours
	CreatedBy   int64        `json:"created_by,omitempty"  db:"created_by"`
	ReporterID  int64        `json:"reporter_id,omitempty" db:"reporter_id"`
	CompletedAt *time.Time   `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time    `json:"created_at,omitempty"  db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"  db:"updated_at"`
}

// TaskComment represents a comment on a task
type TaskComment struct {
	ID        int64     `json:"id,omitempty"         db:"id"`
	TaskID    int64     `json:"task_id,omitempty"    db:"task_id"`
	AuthorID  int64     `json:"author_id,omitempty"  db:"author_id"`
	ParentID  int64     `json:"parent_id,omitempty"  db:"parent_id"`
	Content   string    `json:"content,omitempty"    db:"content"`
	IsEdited  bool      `json:"is_edited,omitempty"  db:"is_edited"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// TaskActivityLog records changes to a task
type TaskActivityLog struct {
	ID        int64           `json:"id,omitempty"         db:"id"`
	TaskID    int64           `json:"task_id,omitempty"    db:"task_id"`
	ActorID   int64           `json:"actor_id,omitempty"   db:"actor_id"`
	Action    string          `json:"action,omitempty"     db:"action"`
	OldValue  json.RawMessage `json:"old_value,omitempty"  db:"old_value"`
	NewValue  json.RawMessage `json:"new_value,omitempty"  db:"new_value"`
	CreatedAt time.Time       `json:"created_at,omitempty" db:"created_at"`
}
