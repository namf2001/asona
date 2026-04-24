package middleware

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

// TokenCheckMiddleware validates the Authorization Bearer token via the auth Controller.
// The controller handles both JWT signature verification and session DB lookup,
// ensuring revoked tokens (e.g. after logout) are rejected immediately.
// Middleware receives auth.Controller — never the repository directly — to
// preserve the Handler → Controller → Repository layering convention.
func TokenCheckMiddleware(ctrl authctrl.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.ERROR.Printf("[TokenCheckMiddleware] missing authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse(
				constants.MissingAuthorizationHeader.Code,
				constants.MissingAuthorizationHeader.Message,
				nil,
			))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			logger.ERROR.Printf("[TokenCheckMiddleware] invalid authorization header: %s", authHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse(
				constants.InvalidAuthorizationHeader.Code,
				constants.InvalidAuthorizationHeader.Message,
				nil,
			))
			return
		}

		tokenStr := headerParts[1]

		userID, email, err := ctrl.ValidateToken(c.Request.Context(), tokenStr)
		if err != nil {
			logger.ERROR.Printf("[TokenCheckMiddleware] token validation failed: %+v", err)
			if errors.Is(err, authctrl.ErrInvalidToken) || errors.Is(err, authctrl.ErrSessionNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse(
					constants.InvalidToken.Code,
					constants.InvalidToken.Message,
					nil,
				))
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewResponse(
				constants.InternalServerError.Code,
				constants.InternalServerError.Message,
				nil,
			))
			return
		}

		c.Set("userID", userID)
		c.Set("email", email)
		c.Next()
	}
}
