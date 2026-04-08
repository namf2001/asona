package chat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	ctrlChannel "asona/internal/controller/chat/channel"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type CreateChannelRequest struct {
	WorkplaceID int64  `json:"workplace_id" binding:"required"`
	ProjectID   int64  `json:"project_id"`
	Name        string `json:"name"         binding:"required"`
	Type        string `json:"type"         binding:"required"`
}

type ChannelResponse struct {
	ID          int64  `json:"id"`
	WorkplaceID int64  `json:"workplace_id"`
	ProjectID   int64  `json:"project_id,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	CreatedBy   int64  `json:"created_by"`
}

// CreateChannel handles POST /api/v1/channels
func (h Handler) CreateChannel(c *gin.Context) {
	var req CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[CreateChannel] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			err.Error(),
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[CreateChannel] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	input := ctrlChannel.CreateChannelInput{
		WorkplaceID: req.WorkplaceID,
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Type:        req.Type,
	}

	ch, err := h.channelCtrl.Create(c.Request.Context(), uid, input)
	if err != nil {
		logger.ERROR.Printf("[CreateChannel] create failed for user %d: %+v", uid, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.CreateChannelFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[CreateChannel] channel created successfully: %s (ID: %d)", ch.Name, ch.ID)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.Success.Code,
		"Channel created",
		ChannelResponse{
			ID:          ch.ID,
			WorkplaceID: ch.WorkplaceID,
			ProjectID:   ch.ProjectID,
			Name:        ch.Name,
			Type:        ch.Type,
			CreatedBy:   ch.CreatedBy,
		},
	))
}

// GetChannel handles GET /api/v1/channels/:id
func (h Handler) GetChannel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		logger.ERROR.Printf("[GetChannel] invalid ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid channel ID",
			nil,
		))
		return
	}

	ch, err := h.channelCtrl.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.ERROR.Printf("[GetChannel] retrieve failed for ID %d: %+v", id, err)
		c.JSON(http.StatusNotFound, response.NewResponse(
			constants.ChannelNotFound.Code,
			constants.ChannelNotFound.Message,
			nil,
		))
		return
	}

	logger.INFO.Printf("[GetChannel] channel retrieved successfully: ID %d", id)
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		ChannelResponse{
			ID:          ch.ID,
			WorkplaceID: ch.WorkplaceID,
			ProjectID:   ch.ProjectID,
			Name:        ch.Name,
			Type:        ch.Type,
			CreatedBy:   ch.CreatedBy,
		},
	))
}
