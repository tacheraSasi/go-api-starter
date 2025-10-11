package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/internals/utils"
)

type PermissionHandler struct {
	permissionService *services.PermissionService
}

func NewPermissionHandler(permissionService *services.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// CreatePermission handles POST /permissions
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Resource    string `json:"resource" binding:"required"`
		Action      string `json:"action" binding:"required"`
		Description string `json:"description,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	permission, err := h.permissionService.CreatePermission(req.Name, req.Resource, req.Action, req.Description)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusCreated, permission)
}

// GetPermission handles GET /permissions/:id
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	permission, err := h.permissionService.GetPermission(uint(id))
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, permission)
}

// ListPermissions handles GET /permissions
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	resource := c.Query("resource")

	permissions, err := h.permissionService.ListPermissions(limit, offset, resource)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, permissions)
}

// UpdatePermission handles PUT /permissions/:id
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	var req struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
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

	permission, err := h.permissionService.UpdatePermission(uint(id), name, description)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, permission)
}

// DeletePermission handles DELETE /permissions/:id
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	err = h.permissionService.DeletePermission(uint(id))
	if err != nil {
		utils.APIError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"message": "Permission deleted successfully"})
}

// GetResourceActions handles GET /permissions/resources/:resource/actions
func (h *PermissionHandler) GetResourceActions(c *gin.Context) {
	resource := c.Param("resource")

	actions, err := h.permissionService.GetResourceActions(resource)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{
		"resource": resource,
		"actions":  actions,
	})
}

// GetAllResources handles GET /permissions/resources
func (h *PermissionHandler) GetAllResources(c *gin.Context) {
	resources, err := h.permissionService.GetAllResources()
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.APISuccess(c, http.StatusOK, gin.H{"resources": resources})
}
