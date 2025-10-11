package services

import (
	"time"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
)

type TokenService interface {
	BlacklistToken(token string, expiresAt time.Time) error
	IsTokenBlacklisted(token string) (bool, error)
}

type tokenService struct {
	repo repositories.TokenRepository
}

func NewTokenService(repo repositories.TokenRepository) TokenService {
	return &tokenService{repo: repo}
}

func (s *tokenService) BlacklistToken(token string, expiresAt time.Time) error {
	blacklistedToken := &models.BlacklistedToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return s.repo.BlacklistToken(blacklistedToken)
}

func (s *tokenService) IsTokenBlacklisted(token string) (bool, error) {
	return s.repo.IsTokenBlacklisted(token)
}
