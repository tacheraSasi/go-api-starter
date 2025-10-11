package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

// PermissionMiddleware checks if the user has the required permission
func PermissionMiddleware(userRepo repositories.UserRepository, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by AuthMiddleware)
		userInterface, exists := c.Get("user")
		if !exists {
			utils.APIError(c, http.StatusUnauthorized, "User not found in context")
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			utils.APIError(c, http.StatusUnauthorized, "Invalid user in context")
			c.Abort()
			return
		}

		// Get user with roles and permissions
		userWithRoles, err := userRepo.GetUserByIDWithRoles(strconv.Itoa(int(user.ID)))
		if err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Failed to get user permissions")
			c.Abort()
			return
		}

		// Check permission
		if !userWithRoles.HasPermission(resource, action) {
			utils.APIError(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole middleware checks if the user has one of the required roles
func RequireRole(userRepo repositories.UserRepository, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by AuthMiddleware)
		userInterface, exists := c.Get("user")
		if !exists {
			utils.APIError(c, http.StatusUnauthorized, "User not found in context")
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			utils.APIError(c, http.StatusUnauthorized, "Invalid user in context")
			c.Abort()
			return
		}

		// Get user with roles
		userWithRoles, err := userRepo.GetUserByIDWithRoles(strconv.Itoa(int(user.ID)))
		if err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Failed to get user roles")
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range requiredRoles {
			if userWithRoles.HasRole(requiredRole) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.APIError(c, http.StatusForbidden, "Insufficient role permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnlyMiddleware is a convenience function for admin-only routes
func AdminOnlyMiddleware(userRepo repositories.UserRepository) gin.HandlerFunc {
	return RequireRole(userRepo, models.RoleAdmin)
}

// ModeratorOrAdminMiddleware allows both moderators and admins
func ModeratorOrAdminMiddleware(userRepo repositories.UserRepository) gin.HandlerFunc {
	return RequireRole(userRepo, models.RoleAdmin, models.RoleModerator)
}
