package chat

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	ctrlMessage "asona/internal/controller/chat/message"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type SendMessageRequest struct {
	ChannelID int64  `json:"channel_id" binding:"required"`
	ParentID  int64  `json:"parent_id"`
	Content   string `json:"content"    binding:"required"`
}

type SendMessageResponse struct {
	ID        int64     `json:"id"`
	ChannelID int64     `json:"channel_id"`
	SenderID  int64     `json:"sender_id"`
	ParentID  int64     `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"is_edited"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SendMessage handles POST /api/v1/messages
// @Summary      Send Chat Message
// @Description  Send a new chat message to a specific channel
// @Tags         chat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      SendMessageRequest  true  "Message details"
// @Success      201      {object}  response.Response{data=SendMessageResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /messages [post]
func (h Handler) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[SendMessage] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			err.Error(),
			nil,
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		logger.ERROR.Printf("[SendMessage] userID not found in context")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidToken.Code,
			constants.InvalidToken.Message,
			nil,
		))
		return
	}
	uid, _ := userID.(int64)

	input := ctrlMessage.SendMessageInput{
		ChannelID: req.ChannelID,
		ParentID:  req.ParentID,
		Content:   req.Content,
	}

	msg, err := h.messageCtrl.Send(c.Request.Context(), uid, input)
	if err != nil {
		logger.ERROR.Printf("[SendMessage] send failed for user %d in channel %d: %+v", uid, req.ChannelID, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.CreateMessageFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	logger.INFO.Printf("[SendMessage] message sent successfully: ID %d in channel %d", msg.ID, msg.ChannelID)
	c.JSON(http.StatusCreated, response.NewResponse(
		constants.Success.Code,
		"Message sent",
		SendMessageResponse{
			ID:        msg.ID,
			ChannelID: msg.ChannelID,
			SenderID:  msg.SenderID,
			ParentID:  msg.ParentID,
			Content:   msg.Content,
			IsEdited:  msg.IsEdited,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		},
	))
}
