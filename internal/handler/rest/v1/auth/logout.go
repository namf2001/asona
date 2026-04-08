package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	authctrl "asona/internal/controller/auth"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

// Logout handles user logout by revoking the current session token.
// @Summary      User Logout
// @Description  Revoke current authenticated session token
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} response.Response
// @Failure      401  {object} response.Response
// @Router       /logout [post]
func (h Handler) Logout(c *gin.Context) {
	token, ok := bearerToken(c.GetHeader("Authorization"))
	if !ok {
		logger.ERROR.Printf("[Logout] invalid authorization header")
		c.JSON(http.StatusUnauthorized, response.NewResponse(
			constants.InvalidAuthorizationHeader.Code,
			constants.InvalidAuthorizationHeader.Message,
			nil,
		))
		return
	}

	err := h.ctrl.Logout(c.Request.Context(), token)
	if err != nil {
		logger.ERROR.Printf("[Logout] failed to revoke token: %+v", err)
		if errors.Is(err, authctrl.ErrSessionNotFound) {
			c.JSON(http.StatusUnauthorized, response.NewResponse(
				constants.InvalidToken.Code,
				constants.InvalidToken.Message,
				nil,
			))
			return
		}

		c.JSON(http.StatusInternalServerError, response.NewResponse(
			constants.InternalServerError.Code,
			constants.InternalServerError.Message,
			nil,
		))
		return
	}

	logger.INFO.Printf("[Logout] user logged out successfully")
	c.JSON(http.StatusOK, response.NewResponse(
		constants.LogoutSuccess.Code,
		constants.LogoutSuccess.Message,
		nil,
	))
}

func bearerToken(authHeader string) (string, bool) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", false
	}
	if strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return parts[1], true
}
