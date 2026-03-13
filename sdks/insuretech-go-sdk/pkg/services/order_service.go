package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// OrderService handles order-related API calls
type OrderService struct {
	Client Client
}

// CreateOrder Create a new order from an approved quotation
func (s *OrderService) CreateOrder(ctx context.Context, req *models.OrderCreationRequest) (*models.OrderCreationResponse, error) {
	path := "/v1/orders"
	var result models.OrderCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListOrders List orders with optional filters
func (s *OrderService) ListOrders(ctx context.Context) (*models.OrdersListingResponse, error) {
	path := "/v1/orders"
	var result models.OrdersListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ConfirmPayment Confirm payment for an order (called by payment gateway callback)
func (s *OrderService) ConfirmPayment(ctx context.Context, orderId string, req *models.PaymentConfirmationRequest) (*models.PaymentConfirmationResponse, error) {
	path := "/v1/orders/{order_id}:confirm-payment"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	var result models.PaymentConfirmationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrderStatus Get lightweight order status
func (s *OrderService) GetOrderStatus(ctx context.Context, orderId string) (*models.OrderStatusRetrievalResponse, error) {
	path := "/v1/orders/{order_id}/status"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	var result models.OrderStatusRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelOrder Cancel an order (only allowed before PAID or POLICY_ISSUED)
func (s *OrderService) CancelOrder(ctx context.Context, orderId string, req *models.OrderCancellationRequest) (*models.OrderCancellationResponse, error) {
	path := "/v1/orders/{order_id}:cancel"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	var result models.OrderCancellationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrder Get a single order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderId string) (*models.OrderRetrievalResponse, error) {
	path := "/v1/orders/{order_id}"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	var result models.OrderRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InitiatePayment Initiate payment for a pending order
func (s *OrderService) InitiatePayment(ctx context.Context, orderId string, req *models.OrdersInitiatePaymentRequest) (*models.OrdersInitiatePaymentResponse, error) {
	path := "/v1/orders/{order_id}:pay"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	var result models.OrdersInitiatePaymentResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

