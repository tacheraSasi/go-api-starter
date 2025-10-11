package dtos

// User DTOs
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	IsActive  bool           `json:"is_active"`
	LastLogin *string        `json:"last_login,omitempty"`
	Roles     []RoleResponse `json:"roles,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

// Role DTOs
type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description,omitempty"`
	PermissionIDs []uint `json:"permission_ids,omitempty"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	IsActive    bool                 `json:"is_active"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

// Permission DTOs
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description,omitempty"`
}

type UpdatePermissionRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Permission Check DTOs
type PermissionCheckRequest struct {
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

type PermissionCheckResponse struct {
	HasPermission bool   `json:"has_permission"`
	Resource      string `json:"resource"`
	Action        string `json:"action"`
}

// Role Assignment DTOs
type AssignRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

type AssignPermissionRequest struct {
	PermissionID uint `json:"permission_id" binding:"required"`
}

// List responses
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type PermissionListResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

// Resource and Action responses
type ResourceActionsResponse struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

type AllResourcesResponse struct {
	Resources []string `json:"resources"`
}
