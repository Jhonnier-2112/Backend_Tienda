package repository

import (
	"time"

	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type DiscountPostgresRepository struct {
	db *gorm.DB
}

func NewDiscountPostgresRepository(db *gorm.DB) *DiscountPostgresRepository {
	return &DiscountPostgresRepository{db: db}
}

func (r *DiscountPostgresRepository) Create(discount *domain.Discount) error {
	return r.db.Create(discount).Error
}

func (r *DiscountPostgresRepository) GetAll() ([]domain.Discount, error) {
	var discounts []domain.Discount
	if err := r.db.Find(&discounts).Error; err != nil {
		return nil, err
	}
	return discounts, nil
}

func (r *DiscountPostgresRepository) GetByID(id uint) (*domain.Discount, error) {
	var discount domain.Discount
	if err := r.db.First(&discount, id).Error; err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *DiscountPostgresRepository) GetByCode(code string) (*domain.Discount, error) {
	var discount domain.Discount
	if err := r.db.Where("code = ?", code).First(&discount).Error; err != nil {
		return nil, err
	}
	return &discount, nil
}

// FindActiveForProduct finds the best active, non-expired discount for a given product
// checking first by productID, then falling back to categoryID.
func (r *DiscountPostgresRepository) FindActiveForProduct(productID uint, categoryID uint) (*domain.Discount, error) {
	now := time.Now()
	var discount domain.Discount

	// Try product-specific discount first
	err := r.db.Where(
		"active = true AND starts_at <= ? AND (expires_at IS NULL OR expires_at >= ?) AND product_id = ?",
		now, now, productID,
	).First(&discount).Error

	if err == nil {
		return &discount, nil
	}

	// Fallback to category discount
	err = r.db.Where(
		"active = true AND starts_at <= ? AND (expires_at IS NULL OR expires_at >= ?) AND category_id = ?",
		now, now, categoryID,
	).First(&discount).Error

	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *DiscountPostgresRepository) Update(discount *domain.Discount) error {
	return r.db.Save(discount).Error
}

func (r *DiscountPostgresRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Discount{}, id).Error
}

func (r *DiscountPostgresRepository) GetPublicActiveDiscounts() ([]domain.Discount, error) {
	var discounts []domain.Discount
	now := time.Now()
	err := r.db.Where(
		"active = true AND starts_at <= ? AND (expires_at IS NULL OR expires_at >= ?) AND code IS NULL",
		now, now,
	).Find(&discounts).Error
	return discounts, err
}
