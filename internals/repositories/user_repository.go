package repositories

import "github.com/tachRoutine/invoice-creator-api/internals/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}
