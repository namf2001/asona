package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Set in header and context for downstream use
		c.Set("X-Request-ID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		
		c.Next()
	}
}
