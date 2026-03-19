package ports

import "tienda-backend/internal/core/domain"

type DashboardRepository interface {
	GetTotalProducts() (int64, error)
	GetTotalOrders() (int64, error)
	GetTotalRevenue() (float64, error)
	GetRecentOrders(limit int) ([]domain.Order, error)
	GetLowStockProducts() ([]domain.Product, error)
	GetActiveUsers() (int64, error)
	GetTopSellingProducts(limit int) ([]domain.ProductSales, error)
	GetAllOrdersForExport() ([]domain.Order, error)
}

type DashboardService interface {
	GetStats() (*domain.DashboardStats, error)
	ExportSalesCSV() ([]byte, error)
}
