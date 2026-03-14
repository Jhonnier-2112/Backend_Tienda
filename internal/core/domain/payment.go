package domain

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrderID        uint           `gorm:"not null;index" json:"order_id"`
	Amount         float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Provider       string         `gorm:"size:50;not null" json:"provider"` // mercadopago, paypal, pse
	Status         PaymentStatus  `gorm:"size:50;not null" json:"status"`
	ProviderID     string         `gorm:"size:255" json:"provider_id"`     // ID from the payment gateway
	ExternalStatus string         `gorm:"size:100" json:"external_status"` // Raw status from provider
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
