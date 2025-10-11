package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole handles POST /roles
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description,omitempty"`
		PermissionIDs []uint `json:"permission_ids,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, err := h.roleService.CreateRole(req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusCreated, role)
}

// GetRole handles GET /roles/:id
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleService.GetRole(uint(id))
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, role)
}

// ListRoles handles GET /roles
func (h *RoleHandler) ListRoles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	activeOnly := c.DefaultQuery("active", "true") == "true"

	roles, err := h.roleService.ListRoles(limit, offset, activeOnly)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, roles)
}

// UpdateRole handles PUT /roles/:id
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		IsActive    *bool   `json:"is_active,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	name := ""
	description := ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}

	role, err := h.roleService.UpdateRole(uint(id), name, description, req.IsActive)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, role)
}

// DeleteRole handles DELETE /roles/:id
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// AddPermissionToRole handles POST /roles/:id/permissions/:permissionId
func (h *RoleHandler) AddPermissionToRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	permissionIDStr := c.Param("permissionId")

	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	err = h.roleService.AddPermissionToRole(uint(roleID), uint(permissionID))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Permission added to role successfully"})
}

// RemovePermissionFromRole handles DELETE /roles/:id/permissions/:permissionId
func (h *RoleHandler) RemovePermissionFromRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	permissionIDStr := c.Param("permissionId")

	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	err = h.roleService.RemovePermissionFromRole(uint(roleID), uint(permissionID))
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Permission removed from role successfully"})
}
