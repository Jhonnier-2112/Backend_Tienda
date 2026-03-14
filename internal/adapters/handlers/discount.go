package handlers

import (
	"net/http"
	"strconv"

	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type DiscountHandler struct {
	service ports.DiscountService
}

func NewDiscountHandler(service ports.DiscountService) *DiscountHandler {
	return &DiscountHandler{service: service}
}

// CreateDiscount godoc
// @Summary Create a Discount
// @Description Creates a new discount (percentage or fixed amount) for a product or category (Admin only)
// @Tags Discounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body domain.Discount true "Discount data"
// @Success 201 {object} domain.Discount
// @Failure 400 {object} map[string]string
// @Router /discounts [post]
func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var input domain.Discount
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	discount, err := h.service.CreateDiscount(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, discount)
}

// GetDiscounts godoc
// @Summary List all Discounts
// @Tags Discounts
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Discount
// @Router /discounts [get]
func (h *DiscountHandler) GetDiscounts(c *gin.Context) {
	discounts, err := h.service.GetDiscounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch discounts"})
		return
	}
	c.JSON(http.StatusOK, discounts)
}

// GetDiscount godoc
// @Summary Get a Discount by ID
// @Tags Discounts
// @Produce json
// @Security BearerAuth
// @Param id path int true "Discount ID"
// @Success 200 {object} domain.Discount
// @Failure 404 {object} map[string]string
// @Router /discounts/{id} [get]
func (h *DiscountHandler) GetDiscount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount ID"})
		return
	}
	discount, err := h.service.GetDiscountByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
		return
	}
	c.JSON(http.StatusOK, discount)
}

// UpdateDiscount godoc
// @Summary Update a Discount
// @Tags Discounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Discount ID"
// @Param input body domain.Discount true "Updated discount data"
// @Success 200 {object} domain.Discount
// @Router /discounts/{id} [put]
func (h *DiscountHandler) UpdateDiscount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount ID"})
		return
	}
	var input domain.Discount
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	discount, err := h.service.UpdateDiscount(uint(id), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, discount)
}

// DeleteDiscount godoc
// @Summary Delete a Discount
// @Tags Discounts
// @Security BearerAuth
// @Param id path int true "Discount ID"
// @Success 200 {object} map[string]string
// @Router /discounts/{id} [delete]
func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount ID"})
		return
	}
	if err := h.service.DeleteDiscount(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Discount deleted successfully"})
}

// GetDiscountByCode godoc
// @Summary Validate a coupon code
// @Description Returns the discount details for a valid, active, non-expired coupon code (public)
// @Tags Discounts
// @Produce json
// @Param code path string true "Coupon code"
// @Success 200 {object} domain.Discount
// @Failure 404 {object} map[string]string
// @Router /discounts/code/{code} [get]
func (h *DiscountHandler) GetDiscountByCode(c *gin.Context) {
	code := c.Param("code")
	discount, err := h.service.GetDiscountByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, discount)
}
func (h *DiscountHandler) GetPublicActiveDiscounts(c *gin.Context) {
	discounts, err := h.service.GetPublicActiveDiscounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active discounts"})
		return
	}
	c.JSON(http.StatusOK, discounts)
}
