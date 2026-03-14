package domain

import (
	"time"
)

type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	URL       string    `gorm:"size:255;not null" json:"url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}
