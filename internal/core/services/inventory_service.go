package services

import (
	"errors"
	"fmt"
	"mime/multipart"

	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
)

type inventoryService struct {
	repo         ports.InventoryRepository
	movementRepo ports.InventoryMovementRepository
	imageStorage ports.ImageStorageService
	auditSvc     ports.AuditLogService
}

func NewInventoryService(
	repo ports.InventoryRepository,
	movementRepo ports.InventoryMovementRepository,
	imageStorage ports.ImageStorageService,
	auditSvc ports.AuditLogService,
) ports.InventoryService {
	return &inventoryService{
		repo:         repo,
		movementRepo: movementRepo,
		imageStorage: imageStorage,
		auditSvc:     auditSvc,
	}
}

// ─── Category Operations ───────────────────────────────────────────────────

func (s *inventoryService) CreateCategory(name, description string) (*domain.Category, error) {
	category := &domain.Category{Name: name, Description: description}
	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *inventoryService) GetCategories() ([]domain.Category, error) {
	return s.repo.GetCategories()
}

func (s *inventoryService) GetCategoryByID(id uint) (*domain.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *inventoryService) UpdateCategory(id uint, name, description string) (*domain.Category, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	category.Name = name
	category.Description = description
	if err := s.repo.UpdateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *inventoryService) DeleteCategory(id uint) error {
	return s.repo.DeleteCategory(id)
}

// ─── Product Operations ────────────────────────────────────────────────────

func (s *inventoryService) CreateProduct(
	userID uint, // Added userID
	name, description, sku string,
	price, costPrice float64,
	stock, minStock int,
	categoryID uint,
	imageFile *multipart.FileHeader,
) (*domain.Product, error) {
	fmt.Printf("Received CreateProduct request: name=%s, categoryID=%d\n", name, categoryID) // Debug log
	if _, err := s.repo.GetCategoryByID(categoryID); err != nil {
		return nil, errors.New("invalid category ID")
	}

	var imageURL *string
	if imageFile != nil {
		if uploadedURL, err := s.imageStorage.UploadImage(imageFile); err == nil {
			imageURL = &uploadedURL
		}
	}

	product := &domain.Product{
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       price,
		CostPrice:   costPrice,
		Stock:       stock,
		MinStock:    minStock,
		CategoryID:  categoryID,
		ImageURL:    imageURL,
	}

	if err := s.repo.CreateProduct(product); err != nil {
		return nil, err
	}

	// AUDIT: Product creation
	s.auditSvc.LogAction(userID, "CREATE_PRODUCT", "Product", product.ID, "", fmt.Sprintf("Name: %s, Price: %.2f", name, price))

	return product, nil
}

func (s *inventoryService) GetProducts() ([]domain.Product, error) {
	return s.repo.GetProducts()
}

func (s *inventoryService) GetProductByID(id uint) (*domain.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *inventoryService) UpdateProduct(
	userID, id uint, // Added userID
	name, description, sku string,
	price, costPrice float64,
	stock, minStock int,
	categoryID uint,
	imageFile *multipart.FileHeader,
) (*domain.Product, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	oldPrice := product.Price // Stored old price for audit
	oldStock := product.Stock // Stored old stock for audit

	if imageFile != nil {
		if uploadedURL, err := s.imageStorage.UploadImage(imageFile); err == nil {
			if product.ImageURL != nil {
				s.imageStorage.DeleteImage(*product.ImageURL)
			}
			product.ImageURL = &uploadedURL
		}
	}

	product.SKU = sku
	product.Name = name
	product.Description = description
	product.Price = price
	product.CostPrice = costPrice
	product.Stock = stock
	product.MinStock = minStock
	product.CategoryID = categoryID

	if err := s.repo.UpdateProduct(product); err != nil {
		return nil, err
	}

	// AUDIT: Price change
	if oldPrice != price {
		s.auditSvc.LogAction(userID, "UPDATE_PRICE", "Product", product.ID, fmt.Sprintf("%.2f", oldPrice), fmt.Sprintf("%.2f", price))
	}

	// AUDIT: Manual stock change from UpdateProduct (usually for corrections in product edit form)
	if oldStock != stock {
		s.auditSvc.LogAction(userID, "UPDATE_STOCK_MANUAL", "Product", product.ID, fmt.Sprintf("%d", oldStock), fmt.Sprintf("%d", stock))
	}

	return product, nil
}

func (s *inventoryService) DeleteProduct(id uint) error {
	product, err := s.repo.GetProductByID(id)
	if err == nil && product.ImageURL != nil {
		s.imageStorage.DeleteImage(*product.ImageURL)
	}
	return s.repo.DeleteProduct(id)
}

// ─── Stock Management ──────────────────────────────────────────────────────

func (s *inventoryService) AdjustStock(userID, productID uint, movementType domain.MovementType, quantity int, note string) (*domain.InventoryMovement, error) {
	if quantity == 0 {
		return nil, errors.New("quantity must be non-zero")
	}

	// Determine delta direction based on type
	delta := quantity
	if movementType == domain.MovementTypeExit {
		delta = -quantity
	}

	// Check resulting stock won't be negative
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	oldStock := product.Stock // Stored old stock for audit
	if oldStock+delta < 0 {
		return nil, errors.New("insufficient stock")
	}

	if err := s.repo.AdjustStock(productID, delta); err != nil {
		return nil, err
	}

	movement := &domain.InventoryMovement{
		ProductID: productID,
		Type:      movementType,
		Quantity:  delta,
		Note:      note,
	}
	if err := s.movementRepo.Create(movement); err != nil {
		return nil, err
	}

	// AUDIT: Stock adjustment
	s.auditSvc.LogAction(userID, "ADJUST_STOCK", "Inventory", productID, fmt.Sprintf("%d", oldStock), fmt.Sprintf("%d", oldStock+delta))

	return movement, nil
}

func (s *inventoryService) GetLowStockProducts() ([]domain.Product, error) {
	return s.repo.GetProductsWithLowStock()
}

func (s *inventoryService) GetMovements(productID *uint) ([]domain.InventoryMovement, error) {
	if productID != nil {
		return s.movementRepo.GetByProductID(*productID)
	}
	return s.movementRepo.GetAll()
}
