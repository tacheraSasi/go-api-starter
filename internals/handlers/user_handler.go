package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserWithRoles(id)
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, user)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	activeOnly := c.DefaultQuery("active", "true") == "true"

	users, err := h.userService.ListUsers(limit, offset, activeOnly)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, users)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name     *string `json:"name,omitempty"`
		Email    *string `json:"email,omitempty"`
		IsActive *bool   `json:"is_active,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	name := ""
	email := ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.Email != nil {
		email = *req.Email
	}

	user, err := h.userService.UpdateUser(id, name, email, req.IsActive)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, user)
}

// UpdateUserPassword handles PUT /users/:id/password
func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.userService.UpdateUserPassword(id, req.Password)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(id)
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// AddRoleToUser handles POST /users/:id/roles/:roleId
func (h *UserHandler) AddRoleToUser(c *gin.Context) {
	userID := c.Param("id")
	roleIDStr := c.Param("roleId")

	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	err = h.userService.AddRoleToUser(userID, uint(roleID))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Role added to user successfully"})
}

// RemoveRoleFromUser handles DELETE /users/:id/roles/:roleId
func (h *UserHandler) RemoveRoleFromUser(c *gin.Context) {
	userID := c.Param("id")
	roleIDStr := c.Param("roleId")

	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	err = h.userService.RemoveRoleFromUser(userID, uint(roleID))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Role removed from user successfully"})
}

// GetUserRoles handles GET /users/:id/roles
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("id")

	roles, err := h.userService.GetUserRoles(userID)
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, roles)
}

// CheckUserPermission handles GET /users/:id/permissions/:resource/:action
func (h *UserHandler) CheckUserPermission(c *gin.Context) {
	userID := c.Param("id")
	resource := c.Param("resource")
	action := c.Param("action")

	hasPermission, err := h.userService.CheckUserPermission(userID, resource, action)
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{
		"has_permission": hasPermission,
		"resource":       resource,
		"action":         action,
	})
}
