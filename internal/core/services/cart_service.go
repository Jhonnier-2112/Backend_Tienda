package services

import (
	"errors"
	"fmt"
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
)

type cartService struct {
	cartRepo     ports.CartRepository
	inventorySvc ports.InventoryService
}

func NewCartService(cartRepo ports.CartRepository, inventorySvc ports.InventoryService) ports.CartService {
	return &cartService{
		cartRepo:     cartRepo,
		inventorySvc: inventorySvc,
	}
}

func (s *cartService) AddToCart(userID, productID uint, quantity int) (*domain.Cart, error) {
	// 1. Validar producto y stock
	product, err := s.inventorySvc.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock: available %d", product.Stock)
	}

	// 2. Obtener carrito del usuario
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		// Asumiremos que el repo maneja la creación si no existe (o lo hacemos aquí)
		return nil, err
	}

	// 3. Revisar si el item ya existe
	var existingItem *domain.CartItem
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			existingItem = &cart.Items[i]
			break
		}
	}

	if existingItem != nil {
		newQuantity := existingItem.Quantity + quantity
		if product.Stock < newQuantity {
			return nil, fmt.Errorf("insufficient total stock: available %d, currently in cart %d", product.Stock, existingItem.Quantity)
		}
		if err := s.cartRepo.UpdateItem(cart.ID, productID, newQuantity); err != nil {
			return nil, err
		}
	} else {
		newItem := &domain.CartItem{
			ProductID: productID,
			Quantity:  quantity,
		}
		if err := s.cartRepo.AddItem(cart.ID, newItem); err != nil {
			return nil, err
		}
	}

	return s.cartRepo.GetByUserID(userID)
}

func (s *cartService) GetMyCart(userID uint) (*domain.Cart, error) {
	return s.cartRepo.GetByUserID(userID)
}

func (s *cartService) RemoveFromCart(userID, productID uint) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if err := s.cartRepo.RemoveItem(cart.ID, productID); err != nil {
		return nil, err
	}

	return s.cartRepo.GetByUserID(userID)
}

func (s *cartService) ClearCart(userID uint) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return err
	}
	return s.cartRepo.Clear(cart.ID)
}
