package repositories

import (
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// Create creates a new permission
func (p *PermissionRepository) Create(permission *models.Permission) error {
	return p.db.Create(permission).Error
}

// GetByID retrieves a permission by ID
func (p *PermissionRepository) GetByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := p.db.First(&permission, id).Error
	return &permission, err
}

// GetByName retrieves a permission by name
func (p *PermissionRepository) GetByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := p.db.Where("name = ?", name).First(&permission).Error
	return &permission, err
}

// GetByResourceAndAction retrieves a permission by resource and action
func (p *PermissionRepository) GetByResourceAndAction(resource, action string) (*models.Permission, error) {
	var permission models.Permission
	err := p.db.Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	return &permission, err
}

// List retrieves all permissions with optional filters
func (p *PermissionRepository) List(limit, offset int, resource string) ([]models.Permission, error) {
	query := p.db.Model(&models.Permission{})

	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	var permissions []models.Permission
	err := query.Limit(limit).Offset(offset).Find(&permissions).Error
	return permissions, err
}

// Update updates a permission
func (p *PermissionRepository) Update(permission *models.Permission) error {
	return p.db.Save(permission).Error
}

// Delete soft deletes a permission
func (p *PermissionRepository) Delete(id uint) error {
	return p.db.Delete(&models.Permission{}, id).Error
}

// GetPermissionsByRole gets all permissions for a specific role
func (p *PermissionRepository) GetPermissionsByRole(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := p.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// GetResourceActions gets all available actions for a resource
func (p *PermissionRepository) GetResourceActions(resource string) ([]string, error) {
	var actions []string
	err := p.db.Model(&models.Permission{}).
		Select("DISTINCT action").
		Where("resource = ?", resource).
		Pluck("action", &actions).Error
	return actions, err
}

// GetAllResources gets all available resources
func (p *PermissionRepository) GetAllResources() ([]string, error) {
	var resources []string
	err := p.db.Model(&models.Permission{}).
		Select("DISTINCT resource").
		Pluck("resource", &resources).Error
	return resources, err
}
