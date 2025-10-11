package services

import (
	"net/http"
	"time"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
)

type AuthService interface {
	Login(email, password string) (models.User, error)
	Register(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	Logout(token string, expiresAt time.Time) error
}

type authService struct {
	repo         repositories.UserRepository
	tokenService TokenService
}

func NewAuthService(repo repositories.UserRepository, tokenService TokenService) AuthService {
	return &authService{repo: repo, tokenService: tokenService}
}

func (s *authService) Login(email, password string) (models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
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
	if err := user.HashPassword(); err != nil {
		return err
	}
	return s.repo.CreateUser(user)
}

func (s *authService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *authService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *authService) Logout(token string, expiresAt time.Time) error {
	return s.tokenService.BlacklistToken(token, expiresAt)
}
