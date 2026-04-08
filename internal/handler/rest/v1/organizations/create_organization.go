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

type CreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

type OrganizationResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Create handles POST /api/v1/organizations
func (h Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[Create] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[Create] userID not found in context")
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
		logger.ERROR.Printf("[Create] create failed for user %d: %+v", uid, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.CreateOrganizationFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[Create] organization created successfully: %s (ID: %d)", org.Name, org.ID)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.Success.Code,
		"Organization created",
		OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			LogoURL:     org.LogoURL,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		},
	))
}
