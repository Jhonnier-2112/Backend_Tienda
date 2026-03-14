package domain

import (
	"time"
)

// Review represents a customer's feedback on a specific product.
type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"` // Relación con el usuario que comenta
	Rating    int       `gorm:"not null" json:"rating"`        // Valor de 1 a 5
	Comment   string    `gorm:"type:text" json:"comment"`
	ImageURL  *string   `gorm:"size:255" json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
