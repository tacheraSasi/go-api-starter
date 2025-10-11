package repositories

import (
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByEmailAndValidatePassword(email, password string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
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

// GetUserByEmail finds a user by their email address
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
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
