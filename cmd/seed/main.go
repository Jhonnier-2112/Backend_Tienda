package main

import (
	"fmt"
	"log"

	"tienda-backend/internal/adapters/repository"
	"tienda-backend/internal/config"
	"tienda-backend/internal/core/services"
)

func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatalf("Failed to load config")
	}

	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repository.NewUserPostgresRepository(db)
	userService := services.NewUserService(userRepo)

	admin, err := userService.CreateAdmin("Admin User", "admin@example.com", "admin123")
	if err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}

	fmt.Printf("Successfully created Admin: %s (ID: %d)\n", admin.Email, admin.ID)
}
