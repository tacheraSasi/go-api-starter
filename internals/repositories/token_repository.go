package repositories

import (
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"gorm.io/gorm"
)

type TokenRepository interface {
	BlacklistToken(token *models.BlacklistedToken) error
	IsTokenBlacklisted(token string) (bool, error)
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) BlacklistToken(token *models.BlacklistedToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&models.BlacklistedToken{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
