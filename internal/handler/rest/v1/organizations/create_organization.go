package organizations

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	ctrlOrgs "asona/internal/controller/organizations"
	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

type CreateOrganizationResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrganization handles POST /api/v1/organizations
// @Summary      Create Organization
// @Description  Create a new organization and assign the creator as ADMIN
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateOrganizationRequest  true  "Organization details"
// @Success      201      {object}  response.Response{data=CreateOrganizationResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /organizations [post]
func (h Handler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[CreateOrganization] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[CreateOrganization] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	org, err := h.ctrl.Create(c.Request.Context(), uid, ctrlOrgs.CreateOrganizationInput{
		Name:        req.Name,
		Description: req.Description,
		LogoURL:     req.LogoURL,
	})
	if err != nil {
		logger.ERROR.Printf("[CreateOrganization] create failed for user %d: %+v", uid, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.CreateOrganizationFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[CreateOrganization] organization created successfully: %s (ID: %d)", org.Name, org.ID)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.Success.Code,
		"Organization created",
		CreateOrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			LogoURL:     org.LogoURL,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		},
	))
}
