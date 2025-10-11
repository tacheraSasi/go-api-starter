package services

import (
	"errors"
	"fmt"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"gorm.io/gorm"
)

type RoleService struct {
	roleRepo       *repositories.RoleRepository
	permissionRepo *repositories.PermissionRepository
}

func NewRoleService(roleRepo *repositories.RoleRepository, permissionRepo *repositories.PermissionRepository) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(name, description string, permissionIDs []uint) (*models.Role, error) {
	// Check if role already exists
	_, err := s.roleRepo.GetByName(name)
	if err == nil {
		return nil, errors.New("role already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Get permissions if provided
	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		for _, permID := range permissionIDs {
			perm, err := s.permissionRepo.GetByID(permID)
			if err != nil {
				return nil, fmt.Errorf("permission with ID %d not found", permID)
			}
			permissions = append(permissions, *perm)
		}
	}

	role := &models.Role{
		Name:        name,
		Description: description,
		IsActive:    true,
		Permissions: permissions,
	}

	err = s.roleRepo.Create(role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// GetRole retrieves a role by ID
func (s *RoleService) GetRole(id uint) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return role, nil
}

// GetRoleByName retrieves a role by name
func (s *RoleService) GetRoleByName(name string) (*models.Role, error) {
	role, err := s.roleRepo.GetByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to perform operation: %w", err)
		}
		return nil, fmt.Errorf("failed to perform operation: %w", err)
	}
	return role, nil
}

// ListRoles retrieves all roles
func (s *RoleService) ListRoles(limit, offset int, activeOnly bool) ([]models.Role, error) {
	roles, err := s.roleRepo.List(limit, offset, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to perform operation: %w", err)
	}
	return roles, nil
}

// UpdateRole updates a role
func (s *RoleService) UpdateRole(id uint, name, description string, isActive *bool) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to perform operation: %w", err)
		}
		return nil, fmt.Errorf("failed to perform operation: %w", err)
	}

	if name != "" {
		// Check if another role with this name exists
		existingRole, err := s.roleRepo.GetByName(name)
		if err == nil && existingRole.ID != id {
			return nil, errors.New("role name already exists")
		}
		role.Name = name
	}

	if description != "" {
		role.Description = description
	}

	if isActive != nil {
		role.IsActive = *isActive
	}

	err = s.roleRepo.Update(role)
	if err != nil {
		return nil, fmt.Errorf("failed to perform operation: %w", err)
	}

	return role, nil
}

// DeleteRole soft deletes a role
func (s *RoleService) DeleteRole(id uint) error {
	_, err := s.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to perform operation: %w", err)
		}
		return fmt.Errorf("failed to perform operation: %w", err)
	}

	err = s.roleRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to perform operation: %w", err)
	}

	return nil
}

// AddPermissionToRole adds a permission to a role
func (s *RoleService) AddPermissionToRole(roleID, permissionID uint) error {
	// Verify role exists
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to perform operation: %w", err)
		}
		return fmt.Errorf("failed to perform operation: %w", err)
	}

	// Verify permission exists
	_, err = s.permissionRepo.GetByID(permissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to perform operation: %w", err)
		}
		return fmt.Errorf("failed to perform operation: %w", err)
	}

	err = s.roleRepo.AddPermission(roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to perform operation: %w", err)
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (s *RoleService) RemovePermissionFromRole(roleID, permissionID uint) error {
	// Verify role exists
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to perform operation: %w", err)
		}
		return fmt.Errorf("failed to perform operation: %w", err)
	}

	err = s.roleRepo.RemovePermission(roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to perform operation: %w", err)
	}

	return nil
}

// InitializeDefaultRoles creates default roles if they don't exist
func (s *RoleService) InitializeDefaultRoles() error {
	defaultRoles := []struct {
		name        string
		description string
	}{
		{models.RoleAdmin, "Administrator with full access"},
		{models.RoleUser, "Regular user with basic access"},
		{models.RoleModerator, "Moderator with limited admin access"},
		{models.RoleGuest, "Guest user with read-only access"},
	}

	for _, roleData := range defaultRoles {
		_, err := s.roleRepo.GetByName(roleData.name)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			role := &models.Role{
				Name:        roleData.name,
				Description: roleData.description,
				IsActive:    true,
			}
			if err := s.roleRepo.Create(role); err != nil {
				return fmt.Errorf("failed to create default role %s: %w", roleData.name, err)
			}
		}
	}

	return nil
}
