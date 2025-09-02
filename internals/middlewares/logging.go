package middlewares

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tachRoutine/invoice-creator-api/pkg/styles"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(styles.Request.Render("Request: %s %s", c.Request.Method, c.Request.URL.Path))
		c.Next()
		log.Println(styles.Response.Render("Response: %d", c.Writer.Status()))
	}
}