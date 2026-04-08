package chat

import (
	"net/http"

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

type CreateChannelResponse struct {
	ID          int64  `json:"id"`
	WorkplaceID int64  `json:"workplace_id"`
	ProjectID   int64  `json:"project_id,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	CreatedBy   int64  `json:"created_by"`
}

// CreateChannel handles POST /api/v1/channels
// @Summary      Create Chat Channel
// @Description  Create a new chat channel within a workplace or project
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateChannelRequest  true  "Channel details"
// @Success      201      {object}  response.Response{data=CreateChannelResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /channels [post]
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
		CreateChannelResponse{
			ID:          ch.ID,
			WorkplaceID: ch.WorkplaceID,
			ProjectID:   ch.ProjectID,
			Name:        ch.Name,
			Type:        ch.Type,
			CreatedBy:   ch.CreatedBy,
		},
	))
}
