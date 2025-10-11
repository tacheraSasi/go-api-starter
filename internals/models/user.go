package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`

	// Legacy field for backward compatibility
	Role string `gorm:"type:varchar(20);default:'user'" json:"role"`
}

func (u *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// HasRole checks if user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName && role.IsActive {
			return true
		}
	}
	return false
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(resource, action string) bool {
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		for _, permission := range role.Permissions {
			if permission.Resource == resource &&
				(permission.Action == action || permission.Action == ActionManage) {
				return true
			}
		}
	}
	return false
}

// GetPermissions returns all permissions for the user
func (u *User) GetPermissions() []Permission {
	permissionMap := make(map[uint]Permission)

	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		for _, permission := range role.Permissions {
			permissionMap[permission.ID] = permission
		}
	}

	permissions := make([]Permission, 0, len(permissionMap))
	for _, permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	return permissions
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// UpdateLastLogin updates the user's last login time
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}
