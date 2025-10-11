package services

import (
	"errors"
	"fmt"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"gorm.io/gorm"
)

type PermissionService struct {
	permissionRepo *repositories.PermissionRepository
}

func NewPermissionService(permissionRepo *repositories.PermissionRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}

// CreatePermission creates a new permission
func (s *PermissionService) CreatePermission(name, resource, action, description string) (*models.Permission, error) {
	// Check if permission already exists
	_, err := s.permissionRepo.GetByResourceAndAction(resource, action)
	if err == nil {
		return nil, errors.New("permission already exists for this resource and action")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing permission: %w", err)
	}

	permission := &models.Permission{
		Name:        name,
		Resource:    resource,
		Action:      action,
		Description: description,
	}

	err = s.permissionRepo.Create(permission)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

// GetPermission retrieves a permission by ID
func (s *PermissionService) GetPermission(id uint) (*models.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	return permission, nil
}

// GetPermissionByName retrieves a permission by name
func (s *PermissionService) GetPermissionByName(name string) (*models.Permission, error) {
	permission, err := s.permissionRepo.GetByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	return permission, nil
}

// ListPermissions retrieves all permissions
func (s *PermissionService) ListPermissions(limit, offset int, resource string) ([]models.Permission, error) {
	permissions, err := s.permissionRepo.List(limit, offset, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	return permissions, nil
}

// UpdatePermission updates a permission
func (s *PermissionService) UpdatePermission(id uint, name, description string) (*models.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	if name != "" {
		permission.Name = name
	}

	if description != "" {
		permission.Description = description
	}

	err = s.permissionRepo.Update(permission)
	if err != nil {
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	return permission, nil
}

// DeletePermission soft deletes a permission
func (s *PermissionService) DeletePermission(id uint) error {
	_, err := s.permissionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return fmt.Errorf("failed to get permission: %w", err)
	}

	err = s.permissionRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}

// InitializeDefaultPermissions creates default permissions if they don't exist
func (s *PermissionService) InitializeDefaultPermissions() error {
	resources := []string{
		models.ResourceUser,
		models.ResourceCustomer,
		models.ResourceInvoice,
		models.ResourceRole,
		models.ResourceSystem,
	}

	actions := []string{
		models.ActionCreate,
		models.ActionRead,
		models.ActionUpdate,
		models.ActionDelete,
		models.ActionList,
		models.ActionManage,
	}

	for _, resource := range resources {
		for _, action := range actions {
			// Skip creating manage permission for system resource except for admin
			if resource == models.ResourceSystem && action != models.ActionManage {
				continue
			}

			name := fmt.Sprintf("%s:%s", resource, action)
			description := fmt.Sprintf("Allow %s on %s", action, resource)

			_, err := s.permissionRepo.GetByResourceAndAction(resource, action)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				permission := &models.Permission{
					Name:        name,
					Resource:    resource,
					Action:      action,
					Description: description,
				}
				if err := s.permissionRepo.Create(permission); err != nil {
					return fmt.Errorf("failed to create default permission %s: %w", name, err)
				}
			}
		}
	}

	return nil
}

// GetResourceActions gets all available actions for a resource
func (s *PermissionService) GetResourceActions(resource string) ([]string, error) {
	actions, err := s.permissionRepo.GetResourceActions(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource actions: %w", err)
	}
	return actions, nil
}

// GetAllResources gets all available resources
func (s *PermissionService) GetAllResources() ([]string, error) {
	resources, err := s.permissionRepo.GetAllResources()
	if err != nil {
		return nil, fmt.Errorf("failed to get all resources: %w", err)
	}
	return resources, nil
}
