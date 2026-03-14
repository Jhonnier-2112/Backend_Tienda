package ports

import (
	"tienda-backend/internal/core/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Update(user *domain.User) error
	FindByVerificationToken(token string) (*domain.User, error)
	FindByResetToken(token string) (*domain.User, error)
}

type InventoryRepository interface {
	// Category Operations
	CreateCategory(category *domain.Category) error
	GetCategories() ([]domain.Category, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	UpdateCategory(category *domain.Category) error
	DeleteCategory(id uint) error

	// Product Operations
	CreateProduct(product *domain.Product) error
	GetProducts() ([]domain.Product, error)
	GetProductByID(id uint) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(id uint) error
	AdjustStock(productID uint, delta int) error
	GetProductsWithLowStock() ([]domain.Product, error)
}

type DiscountRepository interface {
	Create(discount *domain.Discount) error
	GetAll() ([]domain.Discount, error)
	GetByID(id uint) (*domain.Discount, error)
	GetByCode(code string) (*domain.Discount, error)
	FindActiveForProduct(productID uint, categoryID uint) (*domain.Discount, error)
	Update(discount *domain.Discount) error
	Delete(id uint) error
	GetPublicActiveDiscounts() ([]domain.Discount, error)
}

type InventoryMovementRepository interface {
	Create(movement *domain.InventoryMovement) error
	GetByProductID(productID uint) ([]domain.InventoryMovement, error)
	GetAll() ([]domain.InventoryMovement, error)
}

type ReviewRepository interface {
	Create(review *domain.Review) error
	GetByProductID(productID uint) ([]domain.Review, error)
	Delete(id uint) error
}

type OrderRepository interface {
	Create(order *domain.Order) error
	Update(order *domain.Order) error
	GetByID(id uint) (*domain.Order, error)
	GetByUserID(userID uint) ([]domain.Order, error)
	GetAll() ([]domain.Order, error)
}

type CartRepository interface {
	GetByUserID(userID uint) (*domain.Cart, error)
	AddItem(cartID uint, item *domain.CartItem) error
	UpdateItem(cartID, productID uint, quantity int) error
	RemoveItem(cartID, productID uint) error
	Clear(cartID uint) error
}

type AuditLogRepository interface {
	Create(log *domain.AuditLog) error
	GetByEntity(entityName string, entityID uint) ([]domain.AuditLog, error)
	GetAll() ([]domain.AuditLog, error)
}

type RefreshTokenRepository interface {
	Create(token *domain.RefreshToken) error
	FindByToken(token string) (*domain.RefreshToken, error)
	Revoke(token string) error
	RevokeAllForUser(userID uint) error
}

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	Update(payment *domain.Payment) error
	GetByID(id uint) (*domain.Payment, error)
	GetByOrderID(orderID uint) ([]domain.Payment, error)
	GetByProviderID(provider string, providerID string) (*domain.Payment, error)
}

type EmailLogRepository interface {
	Create(log *domain.EmailLog) error
	GetAll() ([]domain.EmailLog, error)
}
