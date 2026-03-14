package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                   uint             `gorm:"primaryKey" json:"id"`
	Name                 string           `gorm:"size:255" json:"-"` // Deprecated, migrating to profiles
	Email                string           `gorm:"size:255;not null;unique" json:"email"`
	Password             string           `gorm:"size:255" json:"-"`
	Role                 string           `gorm:"size:50;not null;default:'customer'" json:"role"` // admin, customer, seller
	IsVerified           bool             `gorm:"default:false" json:"is_verified"`
	VerificationToken    string           `gorm:"size:255" json:"-"`
	ResetPasswordToken   string           `gorm:"size:255" json:"-"`
	ResetPasswordExpires *time.Time       `json:"-"`
	FailedAttempts       int              `gorm:"default:0" json:"-"`
	LockUntil            *time.Time       `json:"-"`
	TwoFactorEnabled     bool             `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret      string           `gorm:"size:255" json:"-"`
	CustomerData         *CustomerProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"customer_profile,omitempty"`
	SellerData           *SellerProfile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"seller_profile,omitempty"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
	DeletedAt            gorm.DeletedAt   `gorm:"index" json:"-"`
}

type CustomerProfile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	FirstName string    `gorm:"size:255;not null" json:"first_name"`
	LastName  string    `gorm:"size:255;not null" json:"last_name"`
	Cedula    string    `gorm:"size:50;not null;unique" json:"cedula"` // DNI/ID
	Phone     string    `gorm:"size:50;default:null" json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SellerProfile struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	CompanyName string    `gorm:"size:255;not null" json:"company_name"`
	ContactName string    `gorm:"size:255;not null" json:"contact_name"`
	NIT         string    `gorm:"size:50;not null;unique" json:"nit"` // Tax ID
	Phone       string    `gorm:"size:50;default:null" json:"phone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null" json:"user_id"`
	Token     string     `gorm:"size:512;not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	UserID       uint             `json:"user_id"`
	AccessToken  string           `json:"access_token,omitempty"`
	RefreshToken string           `json:"refresh_token,omitempty"`
	MfaToken     string           `json:"mfa_token,omitempty"`
	Email        string           `json:"email"`
	Role         string           `json:"role"`
	Customer     *CustomerProfile `json:"customer_profile,omitempty"`
	Seller       *SellerProfile   `json:"seller_profile,omitempty"`
}
