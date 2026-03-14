package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type CartPostgresRepository struct {
	db *gorm.DB
}

func NewCartPostgresRepository(db *gorm.DB) *CartPostgresRepository {
	return &CartPostgresRepository{db: db}
}

func (r *CartPostgresRepository) GetByUserID(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	// Preload items and their products
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new cart if it doesn't exist
			newCart := &domain.Cart{UserID: userID}
			if err := r.db.Create(newCart).Error; err != nil {
				return nil, err
			}
			return newCart, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *CartPostgresRepository) AddItem(cartID uint, item *domain.CartItem) error {
	item.CartID = cartID
	return r.db.Create(item).Error
}

func (r *CartPostgresRepository) UpdateItem(cartID, productID uint, quantity int) error {
	return r.db.Model(&domain.CartItem{}).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Update("quantity", quantity).Error
}

func (r *CartPostgresRepository) RemoveItem(cartID, productID uint) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&domain.CartItem{}).Error
}

func (r *CartPostgresRepository) Clear(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error
}
