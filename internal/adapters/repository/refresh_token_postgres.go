package repository

import (
	"tienda-backend/internal/core/domain"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenPostgresRepository struct {
	db *gorm.DB
}

func NewRefreshTokenPostgresRepository(db *gorm.DB) *RefreshTokenPostgresRepository {
	return &RefreshTokenPostgresRepository{db: db}
}

func (r *RefreshTokenPostgresRepository) Create(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenPostgresRepository) FindByToken(token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	if err := r.db.Where("token = ? AND (revoked_at IS NULL OR revoked_at > ?)", token, time.Now()).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenPostgresRepository) Revoke(token string) error {
	now := time.Now()
	return r.db.Model(&domain.RefreshToken{}).Where("token = ?", token).Update("revoked_at", now).Error
}

func (r *RefreshTokenPostgresRepository) RevokeAllForUser(userID uint) error {
	now := time.Now()
	return r.db.Model(&domain.RefreshToken{}).Where("user_id = ?", userID).Update("revoked_at", now).Error
}
