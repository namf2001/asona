package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/jwt"
)

// TokenCheckMiddleware validates the Authorization Bearer token.
func TokenCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse(
				constants.MissingAuthorizationHeader.Code,
				constants.MissingAuthorizationHeader.Message,
				nil,
			))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse(
				constants.InvalidAuthorizationHeader.Code,
				constants.InvalidAuthorizationHeader.Message,
				nil,
			))
			return
		}

		claims, err := jwt.ParseToken(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewResponse(
				constants.InvalidToken.Code,
				constants.InvalidToken.Message,
				nil,
			))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
