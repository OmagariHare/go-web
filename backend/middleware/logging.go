package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware records detailed information for every HTTP request.
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// After request
		latency := time.Since(start)
		userID := c.GetUint("user_id")
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		if userID > 0 {
			fields = append(fields, zap.Uint("user_id", userID))
		}

		// No need to log c.Errors here, as it's handled by the ErrorHandler middleware.
		// We can, however, log at different levels based on the status code.
		if status >= 500 {
			// Log 5xx server errors with a higher level
			logger.Error("Server Error", fields...)
		} else if status >= 400 {
			// Log 4xx client errors as warnings
			logger.Warn("Client Error", fields...)
		} else {
			// Log successful requests as info
			logger.Info("Request Handled", fields...)
		}
	}
}
