package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/jwt"
	"asona/internal/pkg/logger"
)

// TokenCheckMiddleware validates the Authorization Bearer token.
func TokenCheckMiddleware() gin.HandlerFunc {
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

		claims, err := jwt.ParseToken(headerParts[1])
		if err != nil {
			logger.ERROR.Printf("[TokenCheckMiddleware] parse error for token [%s...]: %+v", headerParts[1][:10], err)
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
