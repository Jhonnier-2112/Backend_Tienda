package ports

import "tienda-backend/internal/core/domain"

type DashboardRepository interface {
	GetTotalProducts() (int64, error)
	GetTotalOrders() (int64, error)
	GetTotalRevenue() (float64, error)
	GetRecentOrders(limit int) ([]domain.Order, error)
	GetLowStockProducts() ([]domain.Product, error)
}

type DashboardService interface {
	GetStats() (*domain.DashboardStats, error)
}
