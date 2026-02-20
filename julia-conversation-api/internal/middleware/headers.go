package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestHeaders propagates mandatory tracking headers
func RequestHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-Id", requestID)
		c.Set("X-Request-Id", requestID)

		correlationID := c.GetHeader("X-Correlation-Id")
		if correlationID == "" {
			correlationID = requestID
		}
		c.Header("X-Correlation-Id", correlationID)
		c.Set("X-Correlation-Id", correlationID)

		c.Next()
	}
}
