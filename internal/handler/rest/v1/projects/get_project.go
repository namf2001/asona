package projects

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type GetProjectResponse struct {
	ID          int64     `json:"id"`
	WorkplaceID int64     `json:"workplace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetProject handles GET /api/v1/projects/:id
// @Summary      Get Project
// @Description  Retrieve details of a specific project by its ID
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  response.Response{data=GetProjectResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /projects/{id} [get]
func (h Handler) GetProject(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		logger.ERROR.Printf("[GetProject] invalid project ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid project ID",
			nil,
		))
		return
	}

	proj, err := h.ctrl.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.ERROR.Printf("[GetProject] retrieval failed for project %d: %+v", id, err)
		c.JSON(http.StatusNotFound, response.NewResponse(
			constants.InternalServerError.Code,
			"Project not found",
			nil,
		))
		return
	}

	logger.INFO.Printf("[GetProject] project retrieved successfully: %d", id)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		GetProjectResponse{
			ID:          proj.ID,
			WorkplaceID: proj.WorkplaceID,
			Name:        proj.Name,
			Description: proj.Description,
			CreatedBy:   proj.CreatedBy,
			CreatedAt:   proj.CreatedAt,
			UpdatedAt:   proj.UpdatedAt,
		},
	))
}
