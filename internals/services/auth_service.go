package services

import (
	"net/http"

	"github.com/tachRoutine/invoice-creator-api/internals/models"
	"github.com/tachRoutine/invoice-creator-api/internals/repositories"
)

type AuthService interface {
	Login(email, password string) (models.User, error)
	Register(user *models.User) error
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(email, password string) (models.User, error){
	user, err := s.repo.GetUserByID(email)
	if err != nil {
		return models.User{}, err
	}
	if err := user.CheckPassword(password); err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (s *authService) Register(user *models.User) error {
	var existingUser *models.User
	existingUser, _ = s.repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return http.ErrBodyNotAllowed
	}
	return nil
}