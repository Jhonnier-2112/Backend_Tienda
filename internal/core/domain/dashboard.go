package domain

type DashboardStats struct {
	TotalProducts   int64         `json:"total_products"`
	TotalOrders     int64         `json:"total_orders"`
	TotalRevenue    float64       `json:"total_revenue"`
	RecentOrders    []Order       `json:"recent_orders"`
	LowStockAlerts  []Product     `json:"low_stock_alerts"`
}
