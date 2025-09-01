package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request details
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		// Call the next handler
		c.Next()

		// Log the response details
		log.Printf("Response: %d", c.Writer.Status())
	}
}