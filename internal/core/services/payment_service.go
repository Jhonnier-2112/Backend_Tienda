package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"tienda-backend/internal/config"
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
	"time"
)

type paymentService struct {
	paymentRepo ports.PaymentRepository
	orderRepo   ports.OrderRepository
	cfg         *config.Config
}

func NewPaymentService(paymentRepo ports.PaymentRepository, orderRepo ports.OrderRepository, cfg *config.Config) ports.PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		cfg:         cfg,
	}
}

func (s *paymentService) CreateMercadoPagoPreference(order *domain.Order) (string, error) {
	if s.cfg.MercadoPagoAccessTok == "" {
		return "", errors.New("mercadopago integration not configured")
	}

	preference := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"title":       fmt.Sprintf("Order #%d", order.ID),
				"description": "Compra en Virtual Store",
				"quantity":    1,
				"currency_id": "COP",
				"unit_price":  order.TotalAmount,
			},
		},
		"external_reference": fmt.Sprintf("%d", order.ID),
		"notification_url":   "https://backend-tienda-wrgv.onrender.com/api/payments/webhook/mercadopago",
		"back_urls": map[string]string{
			"success": s.cfg.AppURL + "/payment/success",
			"pending": s.cfg.AppURL + "/payment/pending",
			"failure": s.cfg.AppURL + "/checkout",
		},
		"auto_return": "approved",
	}

	url := "https://api.mercadopago.com/checkout/preferences"
	payloadBytes, err := json.Marshal(preference)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.cfg.MercadoPagoAccessTok)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Token:", s.cfg.MercadoPagoAccessTok)
	fmt.Println("Payload:", string(payloadBytes))
	fmt.Println("Status:", resp.StatusCode)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("MercadoPago error:", string(bodyBytes))
		return "", fmt.Errorf("mercadopago api error: %s", string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	initPoint, ok := result["init_point"].(string)
	if !ok {
		return "", errors.New("mercadopago response missing init_point")
	}

	providerID, _ := result["id"].(string)

	payment := &domain.Payment{
		OrderID:    order.ID,
		Amount:     order.TotalAmount,
		Provider:   "mercadopago",
		Status:     domain.PaymentStatusPending,
		ProviderID: providerID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return "", err
	}

	return initPoint, nil
}

func (s *paymentService) getPayPalAccessToken() (string, error) {
	url := "https://api-m.sandbox.paypal.com/v1/oauth2/token" // Sandbox URL
	req, err := http.NewRequest("POST", url, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(s.cfg.PayPalClientID, s.cfg.PayPalSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("paypal auth failed with status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func (s *paymentService) CreatePayPalOrder(order *domain.Order) (string, error) {
	if s.cfg.PayPalClientID == "" || s.cfg.PayPalSecret == "" {
		return "", errors.New("paypal integration not configured")
	}

	token, err := s.getPayPalAccessToken()
	if err != nil {
		return "", err
	}

	url := "https://api-m.sandbox.paypal.com/v2/checkout/orders"

	// Ensure amount has 2 decimal places for PayPal USD. If using COP, PayPal might not support it directly,
	// but assuming COP or converted USD value here.
	amountStr := fmt.Sprintf("%.2f", order.TotalAmount)

	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": fmt.Sprintf("%d", order.ID),
				"amount": map[string]interface{}{
					"currency_code": "USD", // Note: PayPal doesn't natively support COP.
					"value":         amountStr,
				},
				"description": fmt.Sprintf("Order #%d from Virtual Store", order.ID),
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("paypal api error: %s", string(bodyBytes))
	}

	var result struct {
		ID    string `json:"id"`
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	var approveLink string
	for _, link := range result.Links {
		if link.Rel == "approve" {
			approveLink = link.Href
			break
		}
	}

	if approveLink == "" {
		return "", errors.New("paypal response missing approve link")
	}

	payment := &domain.Payment{
		OrderID:    order.ID,
		Amount:     order.TotalAmount,
		Provider:   "paypal",
		Status:     domain.PaymentStatusPending,
		ProviderID: result.ID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return "", err
	}

	return approveLink, nil
}

func (s *paymentService) ProcessMercadoPagoWebhook(payload map[string]interface{}) error {
	topic, _ := payload["topic"].(string)
	if topic == "" {
		topic, _ = payload["type"].(string) // Sometimes passed as type
	}
	action, _ := payload["action"].(string)

	// If it's a payment action from Webhook or manual frontend call
	if topic == "payment" || action == "payment.updated" || action == "payment.created" {
		dataMap, ok := payload["data"].(map[string]interface{})
		if !ok {
			return errors.New("invalid payload structure: data missing")
		}
		return s.processPaymentData(dataMap)
	}

	if topic == "merchant_order" {
		// Webhook de merchant_order: hay que consultar la API para obtener los pagos
		resourceURL, _ := payload["resource"].(string)
		if resourceURL == "" {
			return errors.New("merchant_order webhook missing resource URL")
		}

		req, err := http.NewRequest("GET", resourceURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+s.cfg.MercadoPagoAccessTok)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("error fetching merchant_order: %s", string(bodyBytes))
		}

		var orderDetail map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&orderDetail); err != nil {
			return err
		}

		// Procesar todos los pagos de la orden
		payments, ok := orderDetail["payments"].([]interface{})
		if !ok || len(payments) == 0 {
			fmt.Printf("[Webhook MP] merchant_order %v no tiene pagos todavía, se reintentará automáticamente.\n", resourceURL)
			return nil
		}

		for _, p := range payments {
			paymentMap := p.(map[string]interface{})
			if err := s.processPaymentData(paymentMap); err != nil {
				// Logea error pero sigue con otros pagos
				fmt.Printf("[Webhook MP] Error procesando pago %v: %v\n", paymentMap["id"], err)
			}
		}
		return nil
	}

	return fmt.Errorf("unsupported webhook topic/action: %s / %s", topic, action)
}

func (s *paymentService) processPaymentData(dataMap map[string]interface{}) error {
	var mpPaymentID string
	if idStr, ok := dataMap["id"].(string); ok {
		mpPaymentID = idStr
	} else if idFloat, ok := dataMap["id"].(float64); ok {
		mpPaymentID = fmt.Sprintf("%.0f", idFloat)
	}

	extRefStr, _ := dataMap["external_reference"].(string)
	mpStatus, _ := dataMap["status"].(string)

	// If missing external_reference or status, fetch from MercadoPago
	if extRefStr == "" || mpStatus == "" {
		paymentInfo, err := s.getMercadoPagoPayment(mpPaymentID)
		if err == nil && paymentInfo != nil {
			if extRefStr == "" {
				if ext, ok := paymentInfo["external_reference"].(string); ok {
					extRefStr = ext
				}
			}
			if mpStatus == "" {
				if stat, ok := paymentInfo["status"].(string); ok {
					mpStatus = stat
				}
			}
		}
	}

	if mpStatus == "" {
		mpStatus = "approved"
	}

	if extRefStr == "" {
		return fmt.Errorf("missing external_reference for payment %s even after lookup", mpPaymentID)
	}

	var payment domain.Payment
	orderID, err := strconv.Atoi(extRefStr)
	if err != nil {
		return fmt.Errorf("invalid external_reference: %s", extRefStr)
	}
	payments, err := s.paymentRepo.GetByOrderID(uint(orderID))
	if err != nil || len(payments) == 0 {
		return fmt.Errorf("no payments found for order %d", orderID)
	}

	for _, p := range payments {
		if p.Provider == "mercadopago" {
			payment = p
			if p.Status == domain.PaymentStatusPending {
				break
			}
		}
	}
	if payment.ID == 0 {
		return errors.New("mercadopago payment not found for order")
	}

	payment.ExternalStatus = mpStatus
	if mpPaymentID != "" {
		payment.ProviderID = mpPaymentID
	}

	switch mpStatus {
	case "approved":
		payment.Status = domain.PaymentStatusCompleted
	case "rejected", "cancelled":
		payment.Status = domain.PaymentStatusFailed
	case "in_process", "pending":
		payment.Status = domain.PaymentStatusPending
	default:
		payment.Status = domain.PaymentStatusPending
	}

	if err := s.paymentRepo.Update(&payment); err != nil {
		return err
	}

	order, err := s.orderRepo.GetByID(payment.OrderID)
	if err != nil {
		return err
	}

	if payment.Status == domain.PaymentStatusCompleted {
		order.PaymentStatus = domain.PaymentStatusCompleted
		order.Status = domain.OrderStatusPaid
	} else if payment.Status == domain.PaymentStatusFailed {
		order.PaymentStatus = domain.PaymentStatusFailed
	}

	return s.orderRepo.Update(order)
}

func (s *paymentService) ProcessPayPalWebhook(payload map[string]interface{}) error {
	return nil
}

func (s *paymentService) GetOrderPaymentHistory(orderID uint) ([]domain.Payment, error) {
	return s.paymentRepo.GetByOrderID(orderID)
}

func (s *paymentService) getMercadoPagoPayment(paymentID string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.cfg.MercadoPagoAccessTok)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mercadopago verify error: %s", string(body))
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
