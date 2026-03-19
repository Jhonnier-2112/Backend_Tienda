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

func (r *dashboardPostgresRepository) GetActiveUsers() (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("role = ?", "client").Count(&count).Error
	return count, err
}

func (r *dashboardPostgresRepository) GetTopSellingProducts(limit int) ([]domain.ProductSales, error) {
	var sales []domain.ProductSales
	err := r.db.Table("order_items").
		Select("products.name as product_name, products.image_url as product_image, SUM(order_items.quantity) as total_sold, SUM(order_items.quantity * order_items.price) as total_revenue").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.payment_status = ?", domain.PaymentStatusCompleted).
		Group("products.id, products.name, products.image_url").
		Order("total_sold desc").
		Limit(limit).
		Scan(&sales).Error
	return sales, err
}

func (r *dashboardPostgresRepository) GetAllOrdersForExport() ([]domain.Order, error) {
	var orders []domain.Order
	// Preload everything useful for the spreadsheet
	err := r.db.Preload("User").Preload("Items.Product").Order("created_at desc").Find(&orders).Error
	return orders, err
}
