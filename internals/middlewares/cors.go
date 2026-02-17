package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(origins ...string) gin.HandlerFunc {
	allowAll := len(origins) == 0
	allowedOrigins := make(map[string]struct{}, len(origins))

	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "" {
			continue
		}
		if trimmed == "*" {
			allowAll = true
			continue
		}
		allowedOrigins[trimmed] = struct{}{}
	}

	return func(c *gin.Context) {
		requestOrigin := c.GetHeader("Origin")
		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if requestOrigin != "" {
			if _, ok := allowedOrigins[requestOrigin]; ok {
				c.Header("Access-Control-Allow-Origin", requestOrigin)
				c.Header("Vary", "Origin")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
