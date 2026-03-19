package routes

import (
	"tienda-backend/internal/adapters/handlers"
	"tienda-backend/internal/core/ports"
	"tienda-backend/internal/middleware"

	_ "tienda-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	authService ports.AuthService,
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	productHandler *handlers.ProductHandler,
	userHandler *handlers.UserHandler,
	discountHandler *handlers.DiscountHandler,
	reviewHandler *handlers.ReviewHandler,
	orderHandler *handlers.OrderHandler,
	cartHandler *handlers.CartHandler,
	auditHandler *handlers.AuditHandler,
	assistantHandler *handlers.AssistantHandler,
	paymentHandler *handlers.PaymentHandler,
	dashboardHandler *handlers.DashboardHandler,
) *gin.Engine {

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	{
		// ── Public routes ──────────────────────────────────────────
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.GoogleLogin)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/verify", authHandler.VerifyEmail)
			auth.POST("/setup-password", authHandler.SetupPassword)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/verify-2fa", authHandler.Verify2FA)
		}

		// Public read-only
		api.GET("/categories", categoryHandler.GetCategories)
		api.GET("/categories/:id", categoryHandler.GetCategory)
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProduct)
		api.GET("/products/:id/price", productHandler.GetFinalPrice)
		api.GET("/products/:id/reviews", reviewHandler.GetProductReviews)
		api.GET("/discounts/code/:code", discountHandler.GetDiscountByCode)
		api.GET("/discounts/active", discountHandler.GetPublicActiveDiscounts)
		api.POST("/chat/assistant", assistantHandler.HandleMessage)

		// Webhooks (Public)
		api.POST("/payments/webhook/mercadopago", paymentHandler.MercadoPagoWebhook)
		api.POST("/payments/webhook/paypal", paymentHandler.PayPalWebhook)

		// ── Protected routes (valid JWT required) ──────────────────
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Categories (Admin only)
			categoriesAdmin := protected.Group("/categories")
			categoriesAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				categoriesAdmin.POST("", categoryHandler.CreateCategory)
				categoriesAdmin.PUT("/:id", categoryHandler.UpdateCategory)
				categoriesAdmin.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// Dashboard (Admin only)
			dashboardAdmin := protected.Group("/dashboard")
			dashboardAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				dashboardAdmin.GET("/stats", dashboardHandler.GetStats)
			}

			// Products (Admin + Seller can create/update; Admin only delete)
			productsMgmt := protected.Group("/products")
			productsMgmt.Use(middleware.RoleMiddleware("admin", "seller"))
			{
				productsMgmt.POST("", productHandler.CreateProduct)
				productsMgmt.PUT("/:id", productHandler.UpdateProduct)
				productsMgmt.POST("/:id/stock", productHandler.AdjustStock)
				productsMgmt.GET("/:id/movements", productHandler.GetMovements)
			}

			productsAdmin := protected.Group("/products")
			productsAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				productsAdmin.DELETE("/:id", productHandler.DeleteProduct)
				productsAdmin.GET("/low-stock", productHandler.GetLowStock)
			}

			// Reviews (Authenticated users can add)
			protected.POST("/products/:id/reviews", reviewHandler.AddReview)

			// Inventory movements (Admin only)
			inventoryAdmin := protected.Group("/inventory")
			inventoryAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				inventoryAdmin.GET("/movements", productHandler.GetAllMovements)
			}

			// Discounts (Admin only)
			discountsAdmin := protected.Group("/discounts")
			discountsAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				discountsAdmin.POST("", discountHandler.CreateDiscount)
				discountsAdmin.GET("", discountHandler.GetDiscounts)
				discountsAdmin.GET("/:id", discountHandler.GetDiscount)
				discountsAdmin.PUT("/:id", discountHandler.UpdateDiscount)
				discountsAdmin.DELETE("/:id", discountHandler.DeleteDiscount)
			}

			// User Management (Admin only)
			usersAdmin := protected.Group("/users")
			usersAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				usersAdmin.GET("", userHandler.GetAllUsers)
				usersAdmin.POST("/seller", userHandler.CreateSeller)
				usersAdmin.POST("/admin", userHandler.CreateAdmin)
			}

			// Orders
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.PlaceOrder)
				orders.GET("/my", orderHandler.GetMyOrders)
				orders.POST("/:id/pay", orderHandler.InitiatePayment)
				orders.GET("/:id/payments", paymentHandler.GetHistory)
				orders.GET("/:id/receipt/pdf", orderHandler.DownloadReceipt)

				// Admin/Seller Order Management
				ordersAdmin := orders.Group("")
				ordersAdmin.Use(middleware.RoleMiddleware("admin", "seller"))
				{
					ordersAdmin.GET("", orderHandler.GetAllOrders) // Admin mostly, but Seller might need it too
					ordersAdmin.PUT("/:id/status", orderHandler.UpdateStatus)
					ordersAdmin.PUT("/:id/tracking", orderHandler.UpdateShipping)
				}
			}

			// Cart
			cart := protected.Group("/cart")
			{
				cart.GET("", cartHandler.GetMyCart)
				cart.POST("/items", cartHandler.AddToCart)
				cart.DELETE("/items/:product_id", cartHandler.RemoveItem)
				cart.DELETE("", cartHandler.ClearCart)
			}

			// Audit (Admin only)
			audit := protected.Group("/audit")
			audit.Use(middleware.RoleMiddleware("admin"))
			{
				audit.GET("", auditHandler.GetGlobalHistory)
				audit.GET("/products/:id", auditHandler.GetProductHistory)
			}

			// Auth Settings (2FA)
			authSettings := protected.Group("/auth")
			{
				authSettings.POST("/setup-2fa", authHandler.Setup2FA)
				authSettings.POST("/activate-2fa", authHandler.Activate2FA)
			}
		}
	}

	return r
}
