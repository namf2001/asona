package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
)

type ProfileResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// @Summary      Get user profile
// @Description  Get the profile of the currently authenticated user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} response.Response
// @Failure      401  {object} response.Response
// @Router       /profile [get]
func (h Handler) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		ProfileResponse{
			UserID: userID.(string),
			Email:  email.(string),
		},
	))
}
