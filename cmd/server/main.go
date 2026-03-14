package main

import (
	"log"

	"tienda-backend/internal/adapters/handlers"
	"tienda-backend/internal/adapters/repository"
	"tienda-backend/internal/adapters/storage"
	"tienda-backend/internal/config"
	"tienda-backend/internal/core/services"
	"tienda-backend/internal/routes"
)

// @title Virtual Store API
// @version 1.0
// @description Hexagonal Architecture REST API for a Virtual Store with RBAC, Image Uploads, Pricing & Discounts.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 127.0.0.1:10000
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// 1. Configuration
	cfg := config.LoadConfig()

	// 2. Database Connection
	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Fatal error connecting to DB: %v", err)
	}

	// 3. Instantiate Repositories (Driven Adapters)
	userRepo := repository.NewUserPostgresRepository(db)
	inventoryRepo := repository.NewInventoryPostgresRepository(db)
	discountRepo := repository.NewDiscountPostgresRepository(db)
	movementRepo := repository.NewInventoryMovementPostgresRepository(db)
	reviewRepo := repository.NewReviewPostgresRepository(db)
	orderRepo := repository.NewOrderPostgresRepository(db)
	cartRepo := repository.NewCartPostgresRepository(db)
	auditRepo := repository.NewAuditLogPostgresRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenPostgresRepository(db)
	paymentRepo := repository.NewPaymentPostgresRepository(db)
	emailLogRepo := repository.NewEmailLogRepository(db)

	// 4. Instantiate Services (Core Logic)
	imageStorage := storage.NewLocalStorage("uploads", "/uploads")
	emailService := services.NewEmailService(emailLogRepo)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, emailService, cfg.JWTSecret)
	auditService := services.NewAuditLogService(auditRepo)
	inventoryService := services.NewInventoryService(inventoryRepo, movementRepo, imageStorage, auditService)
	discountService := services.NewDiscountService(discountRepo, inventoryRepo)
	reviewService := services.NewReviewService(reviewRepo, inventoryRepo, imageStorage)
	paymentService := services.NewPaymentService(paymentRepo, orderRepo, cfg)
	orderService := services.NewOrderService(orderRepo, inventoryService, discountService)
	cartService := services.NewCartService(cartRepo, inventoryService)
	userService := services.NewUserService(userRepo)
	assistantService := services.NewAssistantService()

	// 5. Instantiate Handlers (Driving Adapters)
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(inventoryService)
	productHandler := handlers.NewProductHandler(inventoryService, discountService)
	userHandler := handlers.NewUserHandler(userService)
	discountHandler := handlers.NewDiscountHandler(discountService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	orderHandler := handlers.NewOrderHandler(orderService, paymentService)
	cartHandler := handlers.NewCartHandler(cartService, discountService)
	auditHandler := handlers.NewAuditHandler(auditService)
	assistantHandler := handlers.NewAssistantHandler(assistantService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// 6. Setup Router
	r := routes.SetupRouter(
		authService,
		authHandler,
		categoryHandler,
		productHandler,
		userHandler,
		discountHandler,
		reviewHandler,
		orderHandler,
		cartHandler,
		auditHandler,
		assistantHandler,
		paymentHandler,
	)

	// Serve Static Files
	// r.Static("/uploads", "./uploads")

	// 7. Start Server
	log.Printf("Server starting on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
