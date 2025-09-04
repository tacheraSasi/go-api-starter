package middlewares

import (
	"github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.WithFields(map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Request received")

		c.Next()

		logger.WithFields(map[string]interface{}{
			"status": c.Writer.Status(),
		}).Info("Response sent")
	}
}
