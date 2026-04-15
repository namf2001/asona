package workplaces

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	ctrlWorkplaces "asona/internal/controller/workplaces"
	"asona/internal/handler/response"
	"asona/internal/model"
	"asona/internal/pkg/logger"
)

type CreateWorkplaceRequest struct {
	Name    string             `json:"name" binding:"required"`
	IconURL string             `json:"icon_url"`
	Size    model.WorkplaceSize `json:"size" binding:"required"`
}

type CreateWorkplaceResponse struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	IconURL   string             `json:"icon_url"`
	Size      model.WorkplaceSize `json:"size"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// CreateWorkplace handles POST /api/v1/workplaces
func (h Handler) CreateWorkplace(c *gin.Context) {
	var req CreateWorkplaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[CreateWorkplace] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[CreateWorkplace] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	wp, err := h.ctrl.Create(c.Request.Context(), uid, ctrlWorkplaces.CreateWorkplaceInput{
		Name:    req.Name,
		IconURL: req.IconURL,
		Size:    req.Size,
	})
	if err != nil {
		logger.ERROR.Printf("[CreateWorkplace] create failed for user %d: %+v", uid, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.CreateWorkplaceFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[CreateWorkplace] workplace created successfully: %s (ID: %d)", wp.Name, wp.ID)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.Success.Code,
		"Workplace created",
		CreateWorkplaceResponse{
			ID:        wp.ID,
			Name:      wp.Name,
			IconURL:   wp.IconURL,
			Size:      wp.Size,
			CreatedAt: wp.CreatedAt,
			UpdatedAt: wp.UpdatedAt,
		},
	))
}
