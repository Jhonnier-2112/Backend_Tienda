package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type paymentPostgresRepository struct {
	db *gorm.DB
}

func NewPaymentPostgresRepository(db *gorm.DB) *paymentPostgresRepository {
	return &paymentPostgresRepository{db: db}
}

func (r *paymentPostgresRepository) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentPostgresRepository) Update(payment *domain.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentPostgresRepository) GetByID(id uint) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentPostgresRepository) GetByOrderID(orderID uint) ([]domain.Payment, error) {
	var payments []domain.Payment
	if err := r.db.Where("order_id = ?", orderID).Order("created_at desc").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentPostgresRepository) GetByProviderID(provider string, providerID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}
