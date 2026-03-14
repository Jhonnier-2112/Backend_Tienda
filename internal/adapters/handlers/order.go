package handlers

import (
	"net/http"
	"strconv"
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service    ports.OrderService
	paymentSvc ports.PaymentService
}

func NewOrderHandler(service ports.OrderService, paymentSvc ports.PaymentService) *OrderHandler {
	return &OrderHandler{
		service:    service,
		paymentSvc: paymentSvc,
	}
}

type PlaceOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	Items           []struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

// PlaceOrder godoc
// @Summary Place a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body PlaceOrderRequest true "Order details"
// @Success 201 {object} domain.Order
// @Router /orders [post]
func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderItems := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		orderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order, err := h.service.PlaceOrder(userID.(uint), req.ShippingAddress, orderItems)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetMyOrders godoc
// @Summary List current user's orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Order
// @Router /orders/my [get]
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, _ := c.Get("userID")

	orders, err := h.service.GetUserOrders(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetAllOrders godoc
// @Summary List all orders (Admin only)
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Order
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateStatus godoc
// @Summary Update order status (Admin only)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param status body map[string]string true "New status"
// @Success 200 {object} domain.Order
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Status domain.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	order, err := h.service.UpdateOrderStatus(userID.(uint), uint(id), domain.OrderStatus(input.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateShipping godoc
// @Summary Update shipping info and set to shipped (Admin/Seller)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} domain.Order
// @Router /orders/{id}/tracking [put]
func (h *OrderHandler) UpdateShipping(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Carrier        string `json:"carrier" binding:"required"`
		TrackingNumber string `json:"tracking_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.UpdateShippingInfo(uint(id), input.Carrier, input.TrackingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// InitiatePayment godoc
// @Summary Initiate payment for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param provider body map[string]string true "Payment provider (mercadopago or paypal)"
// @Success 200 {object} map[string]string "Redirect URL"
// @Router /orders/{id}/pay [post]
func (h *OrderHandler) InitiatePayment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")

	var input struct {
		Provider string `json:"provider" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.GetOrderByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access to order"})
		return
	}

	if order.Status == domain.OrderStatusPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order is already paid"})
		return
	}

	var redirectURL string
	switch input.Provider {
	case "mercadopago":
		redirectURL, err = h.paymentSvc.CreateMercadoPagoPreference(order)
	case "paypal":
		redirectURL, err = h.paymentSvc.CreatePayPalOrder(order)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported payment provider"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirect_url": redirectURL})
}

// DownloadReceipt godoc
// @Summary Download PDF Receipt for an Order
// @Tags Orders
// @Produce application/pdf
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {file} file
// @Router /orders/{id}/receipt/pdf [get]
func (h *OrderHandler) DownloadReceipt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}
	userID, _ := c.Get("userID")

	// Verify order ownership and status
	order, err := h.service.GetOrderByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access to order"})
		return
	}

	if order.Status != domain.OrderStatusPaid && order.Status != domain.OrderStatusShipped && order.Status != domain.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order must be paid to generate a receipt"})
		return
	}

	pdfBytes, err := h.service.GenerateReceiptPDF(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate pdf receipt: " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=receipt-order-"+strconv.Itoa(id)+".pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
