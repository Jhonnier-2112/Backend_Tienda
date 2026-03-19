package repository

import (
	"fmt"
	"log"

	"tienda-backend/internal/config"
	"tienda-backend/internal/core/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-Migrate Models
	err = db.AutoMigrate(
		&domain.User{},
		&domain.CustomerProfile{},
		&domain.SellerProfile{},
		&domain.Category{},
		&domain.ProductImage{},
		&domain.Product{},
		&domain.Discount{},
		&domain.InventoryMovement{},
		&domain.Review{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Cart{},
		&domain.CartItem{},
		&domain.AuditLog{},
		&domain.RefreshToken{},
		&domain.EmailLog{},
		&domain.Payment{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make database migrations: %w", err)
	}

	log.Println("Database connection successfully established")
	return db, nil
}
