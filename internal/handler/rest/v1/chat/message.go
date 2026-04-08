package chat

import (
	"net/http"
	"strconv"
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

type MessageResponse struct {
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
		MessageResponse{
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

// ListMessages handles GET /api/v1/channels/:id/messages
func (h Handler) ListMessages(c *gin.Context) {
	channelID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || channelID == 0 {
		logger.ERROR.Printf("[ListMessages] invalid channel ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid channel ID",
			nil,
		))
		return
	}

	limit := parseIntQuery(c, "limit", 50)
	offset := parseIntQuery(c, "offset", 0)

	msgs, err := h.messageCtrl.List(c.Request.Context(), channelID, limit, offset)
	if err != nil {
		logger.ERROR.Printf("[ListMessages] list failed for channel %d: %+v", channelID, err)
		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.GetMessageFail.Code,
			err.Error(),
			nil,
		))
		return
	}

	result := make([]MessageResponse, 0, len(msgs))
	for _, m := range msgs {
		result = append(result, MessageResponse{
			ID:        m.ID,
			ChannelID: m.ChannelID,
			SenderID:  m.SenderID,
			ParentID:  m.ParentID,
			Content:   m.Content,
			IsEdited:  m.IsEdited,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}
	logger.INFO.Printf("[ListMessages] messages listed successfully for channel %d: count %d", channelID, len(msgs))
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		result,
	))
}

func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	val, err := strconv.Atoi(c.Query(key))
	if err != nil || val < 0 {
		return defaultVal
	}
	return val
}
