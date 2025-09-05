package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		c.Next()
		// End timer
		latency := time.Since(startTime)

		logger.WithFields(map[string]interface{}{
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"client_ip":     c.ClientIP(),
			"user_agent":    c.Request.UserAgent(),
			"content_type":  c.ContentType(),
			"request_id":    c.GetHeader("X-Request-Id"),
			"status":        c.Writer.Status(),
			"latency_ms":    latency.Milliseconds(),
			"response_size": c.Writer.Size(),
		}).Info("Request completed")
	}
}
