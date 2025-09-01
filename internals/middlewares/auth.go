package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tachRoutine/invoice-creator-api/internals/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.APIError(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.APIError(c, http.StatusUnauthorized, "Bearer token is required")
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.APIError(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			utils.APIError(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}