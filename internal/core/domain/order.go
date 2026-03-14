package domain

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pendiente"
	OrderStatusPaid      OrderStatus = "pagado"
	OrderStatusShipped   OrderStatus = "enviado"
	OrderStatusDelivered OrderStatus = "entregado"
	OrderStatusCancelled OrderStatus = "cancelado"
	OrderStatusRefunded  OrderStatus = "reembolsado"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pendiente"
	PaymentStatusProcessing PaymentStatus = "procesando"
	PaymentStatusCompleted  PaymentStatus = "completado"
	PaymentStatusFailed     PaymentStatus = "fallido"
	PaymentStatusRefunded   PaymentStatus = "reembolsado"
)

type Order struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"not null;index" json:"user_id"`
	User             User           `gorm:"foreignKey:UserID" json:"-"`
	Status           OrderStatus    `gorm:"size:50;not null;default:'pending'" json:"status"`
	TotalAmount      float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	ShippingAddress  string         `gorm:"type:text;not null" json:"shipping_address"`
	ShippingCarrier  string         `gorm:"size:100" json:"shipping_carrier,omitempty"`
	TrackingNumber   string         `gorm:"size:100" json:"tracking_number,omitempty"`
	PaymentStatus    PaymentStatus  `gorm:"size:50;not null;default:'pending'" json:"payment_status"`
	PaymentMethod    string         `gorm:"size:50" json:"payment_method,omitempty"`
	PaymentReference string         `gorm:"size:255" json:"payment_reference,omitempty"`
	Items            []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"items"`
	PaymentHistory   []Payment      `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"payment_history,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null;index" json:"order_id"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`     // Price at the time of purchase
	Discount  float64 `gorm:"type:decimal(10,2);default:0" json:"discount"` // Discount at the time of purchase
}
