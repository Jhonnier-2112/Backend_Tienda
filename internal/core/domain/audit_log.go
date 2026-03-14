package domain

import (
	"time"
)

// AuditLog records critical changes made to entities by users.
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Action     string    `gorm:"size:100;not null" json:"action"`      // e.g., "CREATE_PRODUCT", "UPDATE_PRICE", "ADJUST_STOCK"
	EntityName string    `gorm:"size:100;not null" json:"entity_name"` // e.g., "Product", "Inventory"
	EntityID   uint      `gorm:"not null;index" json:"entity_id"`
	OldValue   string    `gorm:"type:text" json:"old_value"`
	NewValue   string    `gorm:"type:text" json:"new_value"`
	CreatedAt  time.Time `json:"created_at"`
}
