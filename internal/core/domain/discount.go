package domain

import (
	"time"

	"gorm.io/gorm"
)

// DiscountType defines how the discount value is applied.
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage" // e.g., 20 = 20% off
	DiscountTypeFixed      DiscountType = "fixed"      // e.g., 5 = $5.00 off
)

// Discount represents a promotional discount applicable to a product or category.
type Discount struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Code       *string        `gorm:"size:100;uniqueIndex" json:"code"` // Optional coupon code
	Type       DiscountType   `gorm:"size:20;not null" json:"type"`
	Value      float64        `gorm:"type:decimal(10,2);not null" json:"value"`
	ProductID  *uint          `gorm:"index" json:"product_id"`  // nil = applies to category
	CategoryID *uint          `gorm:"index" json:"category_id"` // nil = applies to product
	StartsAt   time.Time      `gorm:"not null" json:"starts_at"`
	ExpiresAt  *time.Time     `json:"expires_at"` // nil = no expiration
	Active     bool           `gorm:"default:true" json:"active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
