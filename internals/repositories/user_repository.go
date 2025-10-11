package repositories

import (
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByIDWithRoles(id string) (*models.User, error)
	GetUserByEmailAndValidatePassword(email, password string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByEmailWithRoles(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
	ListUsers(limit, offset int, activeOnly bool) ([]models.User, error)
	AddRoleToUser(userID, roleID uint) error
	RemoveRoleFromUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]models.Role, error)
	UpdateLastLogin(userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser inserts a new user record into the database
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByID finds a user by their unique ID
func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByIDWithRoles finds a user by their unique ID with roles preloaded
func (r *userRepository) GetUserByIDWithRoles(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles.Permissions").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail finds a user by their email address
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmailWithRoles finds a user by their email address with roles preloaded
func (r *userRepository) GetUserByEmailWithRoles(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles.Permissions").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmailAndValidatePassword gets the user data and validate the password returns and error if failed
func (r *userRepository) GetUserByEmailAndValidatePassword(email, password string) (*models.User, error) {
	user, err := r.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := user.CheckPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser saves changes to an existing user
func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// DeleteUser removes a user record from the database by ID
func (r *userRepository) DeleteUser(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// ListUsers retrieves all users with optional filters
func (r *userRepository) ListUsers(limit, offset int, activeOnly bool) ([]models.User, error) {
	query := r.db.Preload("Roles")

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	var users []models.User
	err := query.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// AddRoleToUser adds a role to a user
func (r *userRepository) AddRoleToUser(userID, roleID uint) error {
	return r.db.Model(&models.User{ID: userID}).Association("Roles").Append(&models.Role{ID: roleID})
}

// RemoveRoleFromUser removes a role from a user
func (r *userRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Model(&models.User{ID: userID}).Association("Roles").Delete(&models.Role{ID: roleID})
}

// GetUserRoles gets all roles for a specific user
func (r *userRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").First(&user, userID).Error
	return user.Roles, err
}

// UpdateLastLogin updates the user's last login time
func (r *userRepository) UpdateLastLogin(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login", gorm.Expr("NOW()")).Error
}
