package handlers

import (
	"log"
	"net/http"
	"strconv"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service ports.PaymentService
}

func NewPaymentHandler(service ports.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) MercadoPagoWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	log.Printf("[Webhook MP] payload received: %+v", payload)

	if err := h.service.ProcessMercadoPagoWebhook(payload); err != nil {
		log.Printf("[Webhook MP] ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *PaymentHandler) PayPalWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if err := h.service.ProcessPayPalWebhook(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *PaymentHandler) GetHistory(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))

	history, err := h.service.GetOrderPaymentHistory(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch payment history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
