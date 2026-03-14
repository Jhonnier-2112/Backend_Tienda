package domain

import (
	"time"
)

// Cart represents a user's persistent shopping cart.
type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;uniqueIndex" json:"user_id"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents a specific product and its quantity in a cart.
type CartItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CartID    uint    `gorm:"not null;index" json:"cart_id"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
}
