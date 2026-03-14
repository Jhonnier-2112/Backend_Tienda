package services

import (
	"bytes"
	"errors"
	"fmt"
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"

	"github.com/jung-kurt/gofpdf"
)

type orderService struct {
	orderRepo    ports.OrderRepository
	inventorySvc ports.InventoryService
	discountSvc  ports.DiscountService
}

func NewOrderService(orderRepo ports.OrderRepository, inventorySvc ports.InventoryService, discountSvc ports.DiscountService) ports.OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		inventorySvc: inventorySvc,
		discountSvc:  discountSvc,
	}
}

func (s *orderService) PlaceOrder(userID uint, shippingAddress string, items []domain.OrderItem) (*domain.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("cannot place an order with no items")
	}

	var totalAmount float64
	processedItems := make([]domain.OrderItem, 0, len(items))

	for _, item := range items {
		// 1. Get current product data
		product, err := s.inventorySvc.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}

		// 2. Validate stock
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s (requested %d, available %d)", product.Name, item.Quantity, product.Stock)
		}

		// 3. Get final price with active discounts
		_, discount, _ := s.discountSvc.GetFinalPrice(item.ProductID)

		finalPrice := product.Price
		discountValue := 0.0

		if discount != nil {
			if discount.Type == domain.DiscountTypePercentage {
				discountValue = (product.Price * discount.Value) / 100
			} else {
				discountValue = discount.Value
			}
			finalPrice -= discountValue
		}

		// 4. Update item data with current prices
		item.Price = finalPrice
		item.Discount = discountValue
		totalAmount += finalPrice * float64(item.Quantity)
		processedItems = append(processedItems, item)

		// 5. Adjust stock (Exit)
		_, err = s.inventorySvc.AdjustStock(userID, item.ProductID, domain.MovementTypeExit, item.Quantity, fmt.Sprintf("Order placement by user %d", userID))
		if err != nil {
			return nil, fmt.Errorf("failed to adjust stock for product %d: %w", item.ProductID, err)
		}
	}

	order := &domain.Order{
		UserID:          userID,
		Status:          domain.OrderStatusPending,
		PaymentStatus:   domain.PaymentStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: shippingAddress,
		Items:           processedItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to persist order: %w", err)
	}

	return order, nil
}

func (s *orderService) UpdateOrderStatus(actorUserID, orderID uint, status domain.OrderStatus) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Logic for stock return if cancelled
	if status == domain.OrderStatusCancelled && order.Status != domain.OrderStatusCancelled {
		for _, item := range order.Items {
			_, _ = s.inventorySvc.AdjustStock(actorUserID, item.ProductID, domain.MovementTypeEntry, item.Quantity, fmt.Sprintf("Restored from cancelled order #%d by user %d", orderID, actorUserID))
		}
	}

	order.Status = status
	if status == domain.OrderStatusPaid {
		order.PaymentStatus = domain.PaymentStatusCompleted
	}

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) UpdateShippingInfo(orderID uint, carrier, trackingNumber string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	order.ShippingCarrier = carrier
	order.TrackingNumber = trackingNumber
	order.Status = domain.OrderStatusShipped

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) GetUserOrders(userID uint) ([]domain.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

func (s *orderService) GetAllOrders() ([]domain.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *orderService) GetOrderByID(orderID uint) (*domain.Order, error) {
	return s.orderRepo.GetByID(orderID)
}

func (s *orderService) GenerateReceiptPDF(order *domain.Order) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.CellFormat(190, 10, "COMPROBANTE DE PAGO - VIRTUAL STORE", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Order Info
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(95, 10, fmt.Sprintf("Orden #: %d", order.ID), "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 10, fmt.Sprintf("Fecha: %s", order.CreatedAt.Format("02/01/2006 15:04")), "0", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(190, 10, fmt.Sprintf("Estado: %s", string(order.Status)), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 10, fmt.Sprintf("Cliente (ID): %d", order.UserID), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 10, fmt.Sprintf("Direccion de Envio: %s", order.ShippingAddress), "0", 1, "L", false, 0, "")
	pdf.Ln(10)

	// Items Table Header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(90, 10, "Producto", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Cantidad", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Precio Unit.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Total", "1", 1, "C", true, 0, "")

	// Items Table Body
	pdf.SetFont("Arial", "", 12)
	for _, item := range order.Items {
		productName := "Producto"
		if item.Product.Name != "" {
			productName = item.Product.Name
		} else {
			productName = fmt.Sprintf("Producto ID %d", item.ProductID)
		}

		// Truncate name if too long
		if len(productName) > 35 {
			productName = productName[:32] + "..."
		}

		lineTotal := item.Price * float64(item.Quantity)

		pdf.CellFormat(90, 10, productName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 10, fmt.Sprintf("$%.2f", item.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 10, fmt.Sprintf("$%.2f", lineTotal), "1", 1, "R", false, 0, "")
	}
	pdf.Ln(5)

	// Totals
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(155, 10, "TOTAL PAGADO", "0", 0, "R", false, 0, "")
	pdf.CellFormat(35, 10, fmt.Sprintf("$%.2f", order.TotalAmount), "1", 1, "R", false, 0, "")
	pdf.Ln(15)

	// Payments logic (Optional, since this is a receipt for a paid order, we just show it's paid)
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(190, 10, "Gracias por su compra. Este documento es un comprobante valido de su transaccion.", "0", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
