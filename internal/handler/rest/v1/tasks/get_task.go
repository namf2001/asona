package tasks

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type GetTaskResponse struct {
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

// GetTask handles GET /api/v1/tasks/:id
// @Summary      Get Task
// @Description  Retrieve details of a specific task by its ID
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  response.Response{data=GetTaskResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /tasks/{id} [get]
func (h Handler) GetTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		logger.ERROR.Printf("[GetTask] invalid task ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid task ID",
			nil,
		))
		return
	}

	task, err := h.ctrl.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.ERROR.Printf("[GetTask] retrieval failed for task %d: %+v", id, err)
		c.JSON(http.StatusNotFound, response.NewResponse(
			constants.InternalServerError.Code,
			"Task not found",
			nil,
		))
		return
	}

	logger.INFO.Printf("[GetTask] task retrieved successfully: %d", id)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		GetTaskResponse{
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
