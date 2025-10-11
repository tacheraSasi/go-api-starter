package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
	"github.com/tacheraSasi/go-api-starter/pkg/jwt"
)

func AuthMiddleware(tokenService services.TokenService, jwtSecret []byte) gin.HandlerFunc {
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

		isBlacklisted, err := tokenService.IsTokenBlacklisted(tokenString)
		if err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Failed to check token")
			c.Abort()
			return
		}

		if isBlacklisted {
			utils.APIError(c, http.StatusUnauthorized, "Token is blacklisted")
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			utils.APIError(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		c.Set("userID", claims.User.ID)
		c.Set("userRole", claims.User.Role)
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

