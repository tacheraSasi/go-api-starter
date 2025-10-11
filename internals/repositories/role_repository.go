package repositories

import (
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create creates a new role
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

// GetByName retrieves a role by name
func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	return &role, err
}

// List retrieves all roles with optional filters
func (r *RoleRepository) List(limit, offset int, activeOnly bool) ([]models.Role, error) {
	query := r.db.Preload("Permissions")

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	var roles []models.Role
	err := query.Limit(limit).Offset(offset).Find(&roles).Error
	return roles, err
}

// Update updates a role
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete soft deletes a role
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// AddPermission adds a permission to a role
func (r *RoleRepository) AddPermission(roleID, permissionID uint) error {
	return r.db.Model(&models.Role{ID: roleID}).Association("Permissions").Append(&models.Permission{ID: permissionID})
}

// RemovePermission removes a permission from a role
func (r *RoleRepository) RemovePermission(roleID, permissionID uint) error {
	return r.db.Model(&models.Role{ID: roleID}).Association("Permissions").Delete(&models.Permission{ID: permissionID})
}

// GetUsersWithRole gets all users with a specific role
func (r *RoleRepository) GetUsersWithRole(roleID uint) ([]models.User, error) {
	var role models.Role
	err := r.db.Preload("Users").First(&role, roleID).Error
	return role.Users, err
}
