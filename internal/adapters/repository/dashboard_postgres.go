package repository

import (
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"

	"gorm.io/gorm"
)

type dashboardPostgresRepository struct {
	db *gorm.DB
}

func NewDashboardPostgresRepository(db *gorm.DB) ports.DashboardRepository {
	return &dashboardPostgresRepository{db: db}
}

func (r *dashboardPostgresRepository) GetTotalProducts() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Product{}).Count(&count).Error
	return count, err
}

func (r *dashboardPostgresRepository) GetTotalOrders() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Order{}).Count(&count).Error
	return count, err
}

func (r *dashboardPostgresRepository) GetTotalRevenue() (float64, error) {
	var total struct {
		Sum float64
	}
	// "pagado" or "entregado" could be used, or just payment_status = "completado"
	err := r.db.Model(&domain.Order{}).
		Select("COALESCE(SUM(total_amount), 0) as sum").
		Where("payment_status = ?", domain.PaymentStatusCompleted).
		Scan(&total).Error
	return total.Sum, err
}

func (r *dashboardPostgresRepository) GetRecentOrders(limit int) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Preload("User").Order("created_at desc").Limit(limit).Find(&orders).Error
	return orders, err
}

func (r *dashboardPostgresRepository) GetLowStockProducts() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Where("stock <= min_stock").Find(&products).Error
	return products, err
}
