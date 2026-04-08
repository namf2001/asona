package organizations

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type GetOrganizationResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetOrganization handles GET /api/v1/organizations/:id
// @Summary      Get Organization
// @Description  Retrieve details of a specific organization by its ID
// @Tags         organizations
// @Param        id   path      int  true  "Organization ID"
// @Produce      json
// @Success      200  {object}  response.Response{data=GetOrganizationResponse}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /organizations/{id} [get]
func (h Handler) GetOrganization(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		logger.ERROR.Printf("[GetOrganization] invalid organization ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid organization ID",
			nil,
		))
		return
	}

	org, err := h.ctrl.Get(c.Request.Context(), id)
	if err != nil {
		logger.ERROR.Printf("[Get] retrieval failed for organization %d: %+v", id, err)
		c.JSON(http.StatusNotFound, response.NewResponse(
			constants.InternalServerError.Code,
			"Organization not found",
			nil,
		))
		return
	}

	logger.INFO.Printf("[Get] organization retrieved successfully: %d", id)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		GetOrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			LogoURL:     org.LogoURL,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		},
	))
}
