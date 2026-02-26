package middleware

import (
	"julia-conversation-api/internal/appcontext"

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

		appPlatform := c.GetHeader("X-App-Platform")
		appVersion := c.GetHeader("X-App-Version")

		// Enrich request context with tracking headers
		ctx := c.Request.Context()
		ctx = appcontext.WithHeader(ctx, appcontext.RequestIDKey, requestID)
		ctx = appcontext.WithHeader(ctx, appcontext.CorrelationIDKey, correlationID)
		if appPlatform != "" {
			ctx = appcontext.WithHeader(ctx, appcontext.AppPlatformKey, appPlatform)
		}
		if appVersion != "" {
			ctx = appcontext.WithHeader(ctx, appcontext.AppVersionKey, appVersion)
		}
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
