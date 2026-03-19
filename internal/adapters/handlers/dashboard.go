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

// ExportToExcel returns a CSV file with all sales history
// @Summary Export sales to Excel (CSV)
// @Description Download sales history
// @Tags admin
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {string} string
// @Router /admin/dashboard/export [get]
func (h *DashboardHandler) ExportToExcel(c *gin.Context) {
	data, err := h.service.ExportSalesCSV()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate export file"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=ventas_tienda.csv")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}
