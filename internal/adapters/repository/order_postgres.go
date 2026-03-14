package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type OrderPostgresRepository struct {
	db *gorm.DB
}

func NewOrderPostgresRepository(db *gorm.DB) *OrderPostgresRepository {
	return &OrderPostgresRepository{db: db}
}

func (r *OrderPostgresRepository) Create(order *domain.Order) error {
	// Start transaction to ensure order and items are created together
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *OrderPostgresRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderPostgresRepository) GetByID(id uint) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.Preload("Items.Product").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderPostgresRepository) GetByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderPostgresRepository) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
