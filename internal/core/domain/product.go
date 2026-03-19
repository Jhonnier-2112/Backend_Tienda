package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SKU            string         `gorm:"size:100;uniqueIndex" json:"sku"`
	Name           string         `gorm:"size:255;not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Price          float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	CostPrice      float64        `gorm:"type:decimal(10,2);default:0" json:"cost_price"`
	Stock          int            `gorm:"not null;default:0" json:"stock"`
	MinStock       int            `gorm:"not null;default:5" json:"min_stock"`
	ImageURL       *string        `gorm:"size:255" json:"image_url"`
	CategoryID     uint           `gorm:"not null" json:"category_id"`
	Category       Category       `json:"-"` // Belongs To relationship
	Images         []ProductImage `gorm:"foreignKey:ProductID" json:"images"`
	Specs          string         `gorm:"type:text" json:"specs"` // Stores JSON specifications
	ShippingCost   float64        `gorm:"type:decimal(10,2);default:0" json:"shipping_cost"`
	ShippingOrigin string         `gorm:"size:50;default:'local'" json:"shipping_origin"` // local or international
	Brand          string         `gorm:"size:100" json:"brand"`
	HasPromotion   bool           `gorm:"default:false" json:"has_promotion"`
	IsFreeShipping bool           `gorm:"default:false" json:"is_free_shipping"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductFilter struct {
	Search         string
	CategoryID     *uint
	Brand          string
	HasPromotion   *bool
	IsFreeShipping *bool
	ShippingOrigin string
}
