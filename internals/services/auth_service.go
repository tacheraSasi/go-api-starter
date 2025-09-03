package services

import "github.com/tachRoutine/invoice-creator-api/internals/models"

type AuthService interface {
	Login(email, password string) (models.User, error)
	Register(user *models.User) error
}

