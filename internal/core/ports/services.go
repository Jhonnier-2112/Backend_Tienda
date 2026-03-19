package ports

import (
	"mime/multipart"

	"tienda-backend/internal/core/domain"
)

type AuthService interface {
	RegisterCustomer(firstName, lastName, cedula, email string) (*domain.User, error)
	Login(email, password string) (*domain.LoginResponse, error)
	GoogleLogin(googleToken string) (*domain.LoginResponse, error)
	RefreshToken(token string) (*domain.TokenPair, error)
	ValidateToken(token string) (uint, string, error) // Returns UserID, Role, Error
	VerifyEmail(token string) error
	SetupPassword(token, newPassword string) error
	RequestPasswordReset(email string) error
	ResetPassword(token, newPassword string) error

	// 2FA methods
	Setup2FA(userID uint) (secret string, qrCodeURL string, err error)
	Activate2FA(userID uint, code string) error
	Verify2FA(userID uint, code string) error
	FinalizeLogin(email, password, code string) (*domain.LoginResponse, error)
}

type UserService interface {
	CreateSeller(companyName, contactName, nit, email, password string) (*domain.User, error)
	CreateAdmin(name, email, password string) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
}

type InventoryService interface {
	// Category Operations
	CreateCategory(name, description string) (*domain.Category, error)
	GetCategories() ([]domain.Category, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	UpdateCategory(id uint, name, description string) (*domain.Category, error)
	DeleteCategory(id uint) error

	// Product Operations
	CreateProduct(userID uint, name, description, sku string, price, costPrice float64, stock, minStock int, categoryID uint, brand, shippingOrigin string, shippingCost float64, hasPromotion, isFreeShipping bool, imageFile *multipart.FileHeader) (*domain.Product, error)
	GetProducts(filter *domain.ProductFilter) ([]domain.Product, error)
	GetProductByID(id uint) (*domain.Product, error)
	UpdateProduct(userID, id uint, name, description, sku string, price, costPrice float64, stock, minStock int, categoryID uint, brand, shippingOrigin string, shippingCost float64, hasPromotion, isFreeShipping bool, imageFile *multipart.FileHeader) (*domain.Product, error)
	DeleteProduct(id uint) error

	// Stock Management
	AdjustStock(userID, productID uint, movementType domain.MovementType, quantity int, note string) (*domain.InventoryMovement, error)
	GetLowStockProducts() ([]domain.Product, error)
	GetMovements(productID *uint) ([]domain.InventoryMovement, error)
}

type DiscountService interface {
	CreateDiscount(discount *domain.Discount) (*domain.Discount, error)
	GetDiscounts() ([]domain.Discount, error)
	GetDiscountByID(id uint) (*domain.Discount, error)
	GetDiscountByCode(code string) (*domain.Discount, error)
	UpdateDiscount(id uint, discount *domain.Discount) (*domain.Discount, error)
	DeleteDiscount(id uint) error
	GetFinalPrice(productID uint) (float64, *domain.Discount, error) // price, applied discount, error
	GetPublicActiveDiscounts() ([]domain.Discount, error)
}

type ReviewService interface {
	AddReview(productID, userID uint, rating int, comment string, imageFile *multipart.FileHeader) (*domain.Review, error)
	GetProductReviews(productID uint) ([]domain.Review, error)
}

type OrderService interface {
	PlaceOrder(userID uint, shippingAddress string, items []domain.OrderItem) (*domain.Order, error)
	UpdateOrderStatus(userID, orderID uint, status domain.OrderStatus) (*domain.Order, error)
	UpdateShippingInfo(orderID uint, carrier, trackingNumber string) (*domain.Order, error)
	GetUserOrders(userID uint) ([]domain.Order, error)
	GetAllOrders() ([]domain.Order, error)
	GetOrderByID(orderID uint) (*domain.Order, error)
	GenerateReceiptPDF(order *domain.Order) ([]byte, error)
}

type CartService interface {
	AddToCart(userID, productID uint, quantity int) (*domain.Cart, error)
	GetMyCart(userID uint) (*domain.Cart, error)
	RemoveFromCart(userID, productID uint) (*domain.Cart, error)
	ClearCart(userID uint) error
}

type AuditLogService interface {
	LogAction(userID uint, action, entityName string, entityID uint, oldValue, newValue string) error
	GetProductHistory(productID uint) ([]domain.AuditLog, error)
	GetGlobalHistory() ([]domain.AuditLog, error)
}

type AssistantService interface {
	GetResponse(message string) (string, error)
}

type EmailService interface {
	SendVerificationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
}

type PaymentService interface {
	CreateMercadoPagoPreference(order *domain.Order) (string, error) // Returns redirect URL
	CreatePayPalOrder(order *domain.Order) (string, error)           // Returns redirect URL
	ProcessMercadoPagoWebhook(payload map[string]interface{}) error
	ProcessPayPalWebhook(payload map[string]interface{}) error
	GetOrderPaymentHistory(orderID uint) ([]domain.Payment, error)
}
