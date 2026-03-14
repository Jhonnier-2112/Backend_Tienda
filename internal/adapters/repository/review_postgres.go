package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type ReviewPostgresRepository struct {
	db *gorm.DB
}

func NewReviewPostgresRepository(db *gorm.DB) *ReviewPostgresRepository {
	return &ReviewPostgresRepository{db: db}
}

func (r *ReviewPostgresRepository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewPostgresRepository) GetByProductID(productID uint) ([]domain.Review, error) {
	var reviews []domain.Review
	if err := r.db.Preload("User").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Review{}, id).Error
}
