package domain

type DashboardStats struct {
	TotalProducts   int64          `json:"total_products"`
	TotalOrders     int64          `json:"total_orders"`
	TotalRevenue    float64        `json:"total_revenue"`
	RecentOrders    []Order        `json:"recent_orders"`
	LowStockAlerts  []Product      `json:"low_stock_alerts"`
	ActiveUsers     int64          `json:"active_users"`
	TopProducts     []ProductSales `json:"top_products"`
}

type ProductSales struct {
	ProductName  string  `json:"product_name"`
	ProductImage *string `json:"product_image"`
	TotalSold    int64   `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}
