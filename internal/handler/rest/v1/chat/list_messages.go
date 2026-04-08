package chat

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type ListMessagesResponse struct {
	ID        int64     `json:"id"`
	ChannelID int64     `json:"channel_id"`
	SenderID  int64     `json:"sender_id"`
	ParentID  int64     `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"is_edited"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListMessages handles GET /api/v1/channels/:id/messages
// @Summary      List Chat Messages
// @Description  Retrieve paginated messages for a specific chat channel
// @Tags         chat
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int  true  "Channel ID"
// @Param        limit   query     int  false "Limit (default 50)"
// @Param        offset  query     int  false "Offset (default 0)"
// @Success      200     {object}  response.Response{data=[]ListMessagesResponse}
// @Failure      400     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /channels/{id}/messages [get]
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

	result := make([]ListMessagesResponse, 0, len(msgs))
	for _, m := range msgs {
		result = append(result, ListMessagesResponse{
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
