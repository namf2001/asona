package chat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

type GetChannelResponse struct {
	ID          int64  `json:"id"`
	WorkplaceID int64  `json:"workplace_id"`
	ProjectID   int64  `json:"project_id,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	CreatedBy   int64  `json:"created_by"`
}

// GetChannel handles GET /api/v1/channels/:id
// @Summary      Get Chat Channel
// @Description  Retrieve details of a specific chat channel by its ID
// @Tags         chat
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Channel ID"
// @Success      200  {object}  response.Response{data=GetChannelResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /channels/{id} [get]
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
		GetChannelResponse{
			ID:          ch.ID,
			WorkplaceID: ch.WorkplaceID,
			ProjectID:   ch.ProjectID,
			Name:        ch.Name,
			Type:        ch.Type,
			CreatedBy:   ch.CreatedBy,
		},
	))
}
