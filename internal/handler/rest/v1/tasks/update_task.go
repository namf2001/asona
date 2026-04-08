package tasks

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	ctrlTask "asona/internal/controller/tasks"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeID  int64      `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskResponse struct {
	Message string `json:"message"`
}

// UpdateTask handles PUT /api/v1/tasks/:id
// @Summary      Update Task
// @Description  Update details of an existing task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "Task ID"
// @Param        request  body      UpdateTaskRequest  true  "Updated task details"
// @Success      200      {object}  response.Response{data=UpdateTaskResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /tasks/{id} [put]
func (h Handler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		logger.ERROR.Printf("[UpdateTask] invalid task ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid task ID",
			nil,
		))
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[UpdateTask] failed request param for task %d: %+v", id, err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			err.Error(),
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[UpdateTask] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	err = h.ctrl.Update(c.Request.Context(), uid, id, ctrlTask.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
	})
	if err != nil {
		logger.ERROR.Printf("[UpdateTask] update failed for task %d by user %d: %+v", id, uid, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[UpdateTask] task updated successfully: %d", id)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		"Task updated",
		UpdateTaskResponse{Message: "Task updated successfully"},
	))
}
