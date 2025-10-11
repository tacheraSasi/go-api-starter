package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewUserService(userRepo repositories.UserRepository, roleRepo *repositories.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateUser creates a new user with default role
func (s *UserService) CreateUser(name, email, password string) (*models.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
		IsActive: true,
		Role:     models.RoleUser, // Legacy field
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default role
	defaultRole, err := s.roleRepo.GetByName(models.RoleUser)
	if err == nil {
		_ = s.userRepo.AddRoleToUser(user.ID, defaultRole.ID)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserWithRoles retrieves a user by ID with roles and permissions
func (s *UserService) GetUserWithRoles(id string) (*models.User, error) {
	user, err := s.userRepo.GetUserByIDWithRoles(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByEmailWithRoles retrieves a user by email with roles and permissions
func (s *UserService) GetUserByEmailWithRoles(email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmailWithRoles(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// ListUsers retrieves all users
func (s *UserService) ListUsers(limit, offset int, activeOnly bool) ([]models.User, error) {
	users, err := s.userRepo.ListUsers(limit, offset, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id string, name, email string, isActive *bool) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if name != "" {
		user.Name = name
	}

	if email != "" {
		// Check if another user with this email exists
		existingUser, err := s.userRepo.GetUserByEmail(email)
		if err == nil && existingUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
		user.Email = email
	}

	if isActive != nil {
		user.IsActive = *isActive
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// UpdateUserPassword updates user password
func (s *UserService) UpdateUserPassword(id, newPassword string) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Password = newPassword
	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id string) error {
	_, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.userRepo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// AddRoleToUser adds a role to a user
func (s *UserService) AddRoleToUser(userID string, roleID uint) error {
	// Convert userID string to uint
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Verify user exists
	_, err = s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify role exists
	_, err = s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = s.userRepo.AddRoleToUser(uint(uid), roleID)
	if err != nil {
		return fmt.Errorf("failed to add role to user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *UserService) RemoveRoleFromUser(userID string, roleID uint) error {
	// Convert userID string to uint
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Verify user exists
	_, err = s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.userRepo.RemoveRoleFromUser(uint(uid), roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

// GetUserRoles gets all roles for a user
func (s *UserService) GetUserRoles(userID string) ([]models.Role, error) {
	// Convert userID string to uint
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	roles, err := s.userRepo.GetUserRoles(uint(uid))
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return roles, nil
}

// CheckUserPermission checks if a user has a specific permission
func (s *UserService) CheckUserPermission(userID, resource, action string) (bool, error) {
	user, err := s.GetUserWithRoles(userID)
	if err != nil {
		return false, err
	}

	return user.HasPermission(resource, action), nil
}

// UpdateLastLogin updates the user's last login timestamp
func (s *UserService) UpdateLastLogin(userID string) error {
	// Convert userID string to uint
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user ID")
	}

	err = s.userRepo.UpdateLastLogin(uint(uid))
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}
