package handlers

import (
	"net/http"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service ports.DashboardService
}

func NewDashboardHandler(service ports.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetStats returns aggregated statistics for the admin dashboard.
// @Summary Get admin dashboard statistics
// @Description Fetch aggregated metrics like total sales, orders, and low-stock alerts
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.DashboardStats
// @Router /admin/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
