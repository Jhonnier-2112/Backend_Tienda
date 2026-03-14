package handlers

import (
	"net/http"
	"strconv"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service ports.AuditLogService
}

func NewAuditHandler(service ports.AuditLogService) *AuditHandler {
	return &AuditHandler{service: service}
}

// GetGlobalHistory godoc
// @Summary Get global audit logs
// @Tags Audit
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.AuditLog
// @Router /audit [get]
func (h *AuditHandler) GetGlobalHistory(c *gin.Context) {
	logs, err := h.service.GetGlobalHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetProductHistory godoc
// @Summary Get audit logs for a specific product
// @Tags Audit
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {array} domain.AuditLog
// @Router /audit/products/{id} [get]
func (h *AuditHandler) GetProductHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	logs, err := h.service.GetProductHistory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product history"})
		return
	}
	c.JSON(http.StatusOK, logs)
}
