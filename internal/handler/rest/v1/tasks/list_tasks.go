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

type ListTasksResponse struct {
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

// ListTasks handles GET /api/v1/projects/:id/tasks
// @Summary      List Tasks
// @Description  Retrieve all tasks within a specific project
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  response.Response{data=[]ListTasksResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /projects/{id}/tasks [get]
func (h Handler) ListTasks(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID == 0 {
		logger.ERROR.Printf("[ListTasks] invalid project ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid project ID",
			nil,
		))
		return
	}

	tasks, err := h.ctrl.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		logger.ERROR.Printf("[ListTasks] list failed for project %d: %+v", projectID, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			err.Error(),
			nil,
		))
		return
	}

	result := make([]ListTasksResponse, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, ListTasksResponse{
			ID:          t.ID,
			ProjectID:   t.ProjectID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			Priority:    t.Priority,
			CreatedBy:   t.CreatedBy,
			AssigneeID:  t.AssigneeID,
			DueDate:     t.DueDate,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	logger.INFO.Printf("[ListTasks] tasks listed successfully for project %d: count %d", projectID, len(result))
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		result,
	))
}
