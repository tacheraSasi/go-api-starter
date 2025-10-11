package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null;uniqueIndex" json:"name"`
	Description string         `json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User         `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

type Permission struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null;uniqueIndex" json:"name"`
	Resource    string         `gorm:"not null" json:"resource"`
	Action      string         `gorm:"not null" json:"action"`
	Description string         `json:"description"`
	Roles       []Role         `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	RoleID    uint           `gorm:"not null" json:"role_id"`
	User      User           `json:"user,omitempty"`
	Role      Role           `json:"role,omitempty"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	RoleID       uint           `gorm:"not null" json:"role_id"`
	PermissionID uint           `gorm:"not null" json:"permission_id"`
	Role         Role           `json:"role,omitempty"`
	Permission   Permission     `json:"permission,omitempty"`
}

// Common permission actions
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionList   = "list"
	ActionManage = "manage" // Full access
)

// Common resources
const (
	ResourceUser     = "user"
	ResourceCustomer = "customer"
	ResourceInvoice  = "invoice"
	ResourceRole     = "role"
	ResourceSystem   = "system"
)

// Common roles
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleGuest     = "guest"
)
