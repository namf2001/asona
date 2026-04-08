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

type ListProjectsResponse struct {
	ID          int64     `json:"id"`
	WorkplaceID int64     `json:"workplace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListProjectsByUser handles GET /api/v1/workplaces/:id/projects
// @Summary      List Projects
// @Description  Retrieve all projects within a specific workplace
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Workplace ID"
// @Success      200  {object}  response.Response{data=[]ListProjectsResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /workplaces/{id}/projects [get]
func (h Handler) ListProjects(c *gin.Context) {
	workplaceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workplaceID == 0 {
		logger.ERROR.Printf("[ListProjects] invalid workplace ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid workplace ID",
			nil,
		))
		return
	}

	projs, err := h.ctrl.ListByWorkplace(c.Request.Context(), workplaceID)
	if err != nil {
		logger.ERROR.Printf("[ListProjects] list failed for workplace %d: %+v", workplaceID, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			err.Error(),
			nil,
		))
		return
	}

	result := make([]ListProjectsResponse, 0, len(projs))
	for _, p := range projs {
		result = append(result, ListProjectsResponse{
			ID:          p.ID,
			WorkplaceID: p.WorkplaceID,
			Name:        p.Name,
			Description: p.Description,
			CreatedBy:   p.CreatedBy,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	logger.INFO.Printf("[ListProjects] projects listed successfully for workplace %d: count %d", workplaceID, len(result))
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		result,
	))
}
