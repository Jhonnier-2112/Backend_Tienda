package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type InventoryPostgresRepository struct {
	db *gorm.DB
}

func NewInventoryPostgresRepository(db *gorm.DB) *InventoryPostgresRepository {
	return &InventoryPostgresRepository{db: db}
}

// Category Operations
func (r *InventoryPostgresRepository) CreateCategory(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *InventoryPostgresRepository) GetCategories() ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *InventoryPostgresRepository) GetCategoryByID(id uint) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.Preload("Products").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *InventoryPostgresRepository) UpdateCategory(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *InventoryPostgresRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

// Product Operations
func (r *InventoryPostgresRepository) CreateProduct(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *InventoryPostgresRepository) GetProducts() ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *InventoryPostgresRepository) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *InventoryPostgresRepository) UpdateProduct(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *InventoryPostgresRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

// AdjustStock atomically changes the stock count by delta (positive or negative).
func (r *InventoryPostgresRepository) AdjustStock(productID uint, delta int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock", r.db.Raw("stock + ?", delta)).Error
}

// GetProductsWithLowStock returns products whose stock is at or below their min_stock threshold.
func (r *InventoryPostgresRepository) GetProductsWithLowStock() ([]domain.Product, error) {
	var products []domain.Product
	if err := r.db.Preload("Category").
		Where("stock <= min_stock").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
