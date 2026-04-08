package tasks

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	ctrlTask "asona/internal/controller/tasks"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type CreateTaskRequest struct {
	ProjectID   int64      `json:"project_id"  binding:"required"`
	Title       string     `json:"title"       binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeID  int64      `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

type CreateTaskResponse struct {
	ID          int64      `json:"id"`
	ProjectID   int64      `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	CreatedBy   int64      `json:"created_by"`
	AssigneeID  int64      `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTask handles POST /api/v1/tasks
// @Summary      Create Task
// @Description  Create a new task in a specific project
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateTaskRequest  true  "Task details"
// @Success      201      {object}  response.Response{data=CreateTaskResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /tasks [post]
func (h Handler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[CreateTask] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[CreateTask] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	task, err := h.ctrl.Create(c.Request.Context(), uid, ctrlTask.CreateTaskInput{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
	})
	if err != nil {
		logger.ERROR.Printf("[CreateTask] create failed for user %d in project %d: %+v", uid, req.ProjectID, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[CreateTask] task created successfully: %s (ID: %d)", task.Title, task.ID)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.Success.Code,
		"Task created",
		CreateTaskResponse{
			ID:          task.ID,
			ProjectID:   task.ProjectID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			CreatedBy:   task.CreatedBy,
			AssigneeID:  task.AssigneeID,
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
	))
}
