package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
)

type dashboardService struct {
	repo ports.DashboardRepository
}

func NewDashboardService(repo ports.DashboardRepository) ports.DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetStats() (*domain.DashboardStats, error) {
	productsChan := make(chan int64)
	ordersChan := make(chan int64)
	revenueChan := make(chan float64)
	usersChan := make(chan int64)

	var productsErr, ordersErr, revenueErr, usersErr error

	go func() {
		count, err := s.repo.GetTotalProducts()
		productsErr = err
		productsChan <- count
	}()

	go func() {
		count, err := s.repo.GetTotalOrders()
		ordersErr = err
		ordersChan <- count
	}()

	go func() {
		total, err := s.repo.GetTotalRevenue()
		revenueErr = err
		revenueChan <- total
	}()

	go func() {
		count, err := s.repo.GetActiveUsers()
		usersErr = err
		usersChan <- count
	}()

	totalProducts := <-productsChan
	totalOrders := <-ordersChan
	totalRevenue := <-revenueChan
	activeUsers := <-usersChan

	if productsErr != nil {
		return nil, productsErr
	}
	if ordersErr != nil {
		return nil, ordersErr
	}
	if revenueErr != nil {
		return nil, revenueErr
	}
	if usersErr != nil {
		return nil, usersErr
	}

	recentOrders, err := s.repo.GetRecentOrders(5)
	if err != nil {
		return nil, err
	}

	lowStock, err := s.repo.GetLowStockProducts()
	if err != nil {
		return nil, err
	}

	topProducts, err := s.repo.GetTopSellingProducts(5)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardStats{
		TotalProducts:  totalProducts,
		TotalOrders:    totalOrders,
		TotalRevenue:   totalRevenue,
		RecentOrders:   recentOrders,
		LowStockAlerts: lowStock,
		ActiveUsers:    activeUsers,
		TopProducts:    topProducts,
	}, nil
}

func (s *dashboardService) ExportSalesCSV() ([]byte, error) {
	orders, err := s.repo.GetAllOrdersForExport()
	if err != nil {
		return nil, err
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	// Header
	w.Write([]string{"ID de Orden", "Fecha", "Cliente Email", "Total ($)", "Estado Orden", "Estado Pago", "Método de Pago", "Productos Comprados"})

	for _, order := range orders {
		var productsList string
		for i, item := range order.Items {
			if i > 0 {
				productsList += " | "
			}
			productsList += fmt.Sprintf("%d x %s", item.Quantity, item.Product.Name)
		}

		w.Write([]string{
			fmt.Sprintf("%d", order.ID),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.User.Email,
			fmt.Sprintf("%.2f", order.TotalAmount),
			string(order.Status),
			string(order.PaymentStatus),
			order.PaymentMethod,
			productsList,
		})
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
