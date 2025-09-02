package middlewares

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tachRoutine/invoice-creator-api/pkg/styles"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(styles.Request.Render("Request: %s %s", c.Request.Method, c.Request.URL.Path))
		c.Next()
		log.Println(styles.Response.Render("Response: %s", strconv.Itoa(c.Writer.Status()))) //TODO: Add response body logging
	}
}
