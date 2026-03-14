package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type InventoryMovementPostgresRepository struct {
	db *gorm.DB
}

func NewInventoryMovementPostgresRepository(db *gorm.DB) *InventoryMovementPostgresRepository {
	return &InventoryMovementPostgresRepository{db: db}
}

func (r *InventoryMovementPostgresRepository) Create(movement *domain.InventoryMovement) error {
	return r.db.Create(movement).Error
}

func (r *InventoryMovementPostgresRepository) GetByProductID(productID uint) ([]domain.InventoryMovement, error) {
	var movements []domain.InventoryMovement
	if err := r.db.Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&movements).Error; err != nil {
		return nil, err
	}
	return movements, nil
}

func (r *InventoryMovementPostgresRepository) GetAll() ([]domain.InventoryMovement, error) {
	var movements []domain.InventoryMovement
	if err := r.db.Order("created_at DESC").
		Find(&movements).Error; err != nil {
		return nil, err
	}
	return movements, nil
}
