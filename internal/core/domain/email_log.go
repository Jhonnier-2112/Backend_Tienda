package domain

import (
	"time"

	"gorm.io/gorm"
)

type EmailLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ToAddress string         `gorm:"size:255;not null" json:"to_address"`
	Subject   string         `gorm:"size:255;not null" json:"subject"`
	Body      string         `gorm:"type:text;not null" json:"body"`
	Status    string         `gorm:"size:50;not null" json:"status"` // e.g. "sent", "failed"
	Error     string         `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
