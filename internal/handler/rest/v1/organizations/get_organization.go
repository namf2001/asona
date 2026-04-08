package organizations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
)

// Get handles GET /api/v1/organizations/:id
func (h Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			"Invalid organization ID",
			nil,
		))
		return
	}

	org, err := h.ctrl.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(
			constants.InternalServerError.Code,
			"Organization not found",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		org,
	))
}
