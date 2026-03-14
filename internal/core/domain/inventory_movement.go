package domain

import "time"

// MovementType classifies the reason for a stock change.
type MovementType string

const (
	MovementTypeEntry      MovementType = "entry"      // stock entering (purchase, return)
	MovementTypeExit       MovementType = "exit"       // stock leaving (sale, damage)
	MovementTypeAdjustment MovementType = "adjustment" // manual correction
)

// InventoryMovement records every change in the stock of a product.
type InventoryMovement struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	ProductID uint         `gorm:"not null;index" json:"product_id"`
	Product   Product      `gorm:"foreignKey:ProductID" json:"-"`
	Type      MovementType `gorm:"size:20;not null" json:"type"`
	Quantity  int          `gorm:"not null" json:"quantity"` // positive = entry, negative = exit/adjustment
	Note      string       `gorm:"size:500" json:"note"`
	CreatedAt time.Time    `json:"created_at"`
}
