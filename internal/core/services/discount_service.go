package services

import (
	"errors"
	"time"

	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
)

type discountService struct {
	discountRepo  ports.DiscountRepository
	inventoryRepo ports.InventoryRepository
}

func NewDiscountService(discountRepo ports.DiscountRepository, inventoryRepo ports.InventoryRepository) ports.DiscountService {
	return &discountService{
		discountRepo:  discountRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *discountService) CreateDiscount(discount *domain.Discount) (*domain.Discount, error) {
	if discount.Type != domain.DiscountTypePercentage && discount.Type != domain.DiscountTypeFixed {
		return nil, errors.New("discount type must be 'percentage' or 'fixed'")
	}
	if discount.Value <= 0 {
		return nil, errors.New("discount value must be greater than zero")
	}
	if discount.Type == domain.DiscountTypePercentage && discount.Value > 100 {
		return nil, errors.New("percentage discount cannot exceed 100")
	}
	if discount.ProductID == nil && discount.CategoryID == nil {
		return nil, errors.New("discount must be linked to a product or a category")
	}
	if discount.ProductID != nil && discount.CategoryID != nil {
		return nil, errors.New("discount can only be linked to a product OR a category, not both")
	}

	if err := s.discountRepo.Create(discount); err != nil {
		return nil, err
	}
	return discount, nil
}

func (s *discountService) GetDiscounts() ([]domain.Discount, error) {
	return s.discountRepo.GetAll()
}

func (s *discountService) GetDiscountByID(id uint) (*domain.Discount, error) {
	return s.discountRepo.GetByID(id)
}

func (s *discountService) GetDiscountByCode(code string) (*domain.Discount, error) {
	discount, err := s.discountRepo.GetByCode(code)
	if err != nil {
		return nil, errors.New("coupon code not found")
	}
	now := time.Now()
	if !discount.Active {
		return nil, errors.New("coupon is inactive")
	}
	if now.Before(discount.StartsAt) {
		return nil, errors.New("coupon is not yet valid")
	}
	if discount.ExpiresAt != nil && now.After(*discount.ExpiresAt) {
		return nil, errors.New("coupon has expired")
	}
	return discount, nil
}

func (s *discountService) UpdateDiscount(id uint, update *domain.Discount) (*domain.Discount, error) {
	existing, err := s.discountRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("discount not found")
	}

	existing.Name = update.Name
	existing.Type = update.Type
	existing.Value = update.Value
	existing.StartsAt = update.StartsAt
	existing.ExpiresAt = update.ExpiresAt
	existing.Active = update.Active
	existing.ProductID = update.ProductID
	existing.CategoryID = update.CategoryID
	existing.Code = update.Code

	if err := s.discountRepo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *discountService) DeleteDiscount(id uint) error {
	return s.discountRepo.Delete(id)
}

// GetFinalPrice returns the sale price for the product after applying any active discount.
func (s *discountService) GetFinalPrice(productID uint) (float64, *domain.Discount, error) {
	product, err := s.inventoryRepo.GetProductByID(productID)
	if err != nil {
		return 0, nil, errors.New("product not found")
	}

	discount, err := s.discountRepo.FindActiveForProduct(productID, product.CategoryID)
	if err != nil {
		// No active discount found — return original price
		return product.Price, nil, nil
	}

	finalPrice := applyDiscount(product.Price, discount)
	return finalPrice, discount, nil
}

// applyDiscount calculates the discounted price.
func applyDiscount(originalPrice float64, d *domain.Discount) float64 {
	var finalPrice float64
	switch d.Type {
	case domain.DiscountTypePercentage:
		finalPrice = originalPrice * (1 - d.Value/100)
	case domain.DiscountTypeFixed:
		finalPrice = originalPrice - d.Value
	default:
		finalPrice = originalPrice
	}
	if finalPrice < 0 {
		finalPrice = 0
	}
	return finalPrice
}
func (s *discountService) GetPublicActiveDiscounts() ([]domain.Discount, error) {
	return s.discountRepo.GetPublicActiveDiscounts()
}
