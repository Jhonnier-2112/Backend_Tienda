package handlers

import (
	"net/http"
	"strconv"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service     ports.CartService
	discountSvc ports.DiscountService
}

func NewCartHandler(service ports.CartService, discountSvc ports.DiscountService) *CartHandler {
	return &CartHandler{
		service:     service,
		discountSvc: discountSvc,
	}
}

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// AddToCart godoc
// @Summary Add or update item in cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body AddToCartRequest true "Item details"
// @Success 200 {object} domain.Cart
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.AddToCart(userID.(uint), req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// GetMyCart godoc
// @Summary Get current user's cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /cart [get]
func (h *CartHandler) GetMyCart(c *gin.Context) {
	userID, _ := c.Get("userID")

	cart, err := h.service.GetMyCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	// Calculate prices and total dynamically
	var total float64
	type cartItemResponse struct {
		ProductID       uint    `json:"product_id"`
		ProductName     string  `json:"product_name"`
		ImageURL        *string `json:"image_url"`
		OriginalPrice   float64 `json:"original_price"`
		DiscountedPrice float64 `json:"discounted_price"`
		Quantity        int     `json:"quantity"`
		Subtotal        float64 `json:"subtotal"`
	}

	itemsResponse := make([]cartItemResponse, 0, len(cart.Items))
	for _, item := range cart.Items {
		finalPrice, _, _ := h.discountSvc.GetFinalPrice(item.ProductID)
		subtotal := finalPrice * float64(item.Quantity)
		total += subtotal

		itemsResponse = append(itemsResponse, cartItemResponse{
			ProductID:       item.ProductID,
			ProductName:     item.Product.Name,
			ImageURL:        item.Product.ImageURL,
			OriginalPrice:   item.Product.Price,
			DiscountedPrice: finalPrice,
			Quantity:        item.Quantity,
			Subtotal:        subtotal,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": itemsResponse,
		"total": total,
	})
}

// RemoveItem godoc
// @Summary Remove an item from cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param product_id path int true "Product ID"
// @Success 200 {object} domain.Cart
// @Router /cart/items/{product_id} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	productID, _ := strconv.Atoi(c.Param("product_id"))

	cart, err := h.service.RemoveFromCart(userID.(uint), uint(productID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// ClearCart godoc
// @Summary Empty the cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, _ := c.Get("userID")

	if err := h.service.ClearCart(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.Status(http.StatusNoContent)
}
