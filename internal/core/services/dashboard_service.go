package services

import (
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

	var productsErr, ordersErr, revenueErr error

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

	totalProducts := <-productsChan
	totalOrders := <-ordersChan
	totalRevenue := <-revenueChan

	if productsErr != nil {
		return nil, productsErr
	}
	if ordersErr != nil {
		return nil, ordersErr
	}
	if revenueErr != nil {
		return nil, revenueErr
	}

	recentOrders, err := s.repo.GetRecentOrders(5)
	if err != nil {
		return nil, err
	}

	lowStock, err := s.repo.GetLowStockProducts()
	if err != nil {
		return nil, err
	}

	return &domain.DashboardStats{
		TotalProducts:  totalProducts,
		TotalOrders:    totalOrders,
		TotalRevenue:   totalRevenue,
		RecentOrders:   recentOrders,
		LowStockAlerts: lowStock,
	}, nil
}
