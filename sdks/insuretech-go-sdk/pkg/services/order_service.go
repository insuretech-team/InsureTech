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

// ListOrders List orders with optional filters
func (s *OrderService) ListOrders(ctx context.Context) error {
	path := "/v1/orders"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateOrder Create a new order from an approved quotation
func (s *OrderService) CreateOrder(ctx context.Context, req *models.OrderCreationRequest) error {
	path := "/v1/orders"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetOrder Get a single order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderId string) error {
	path := "/v1/orders/{order_id}"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetOrderStatus Get lightweight order status
func (s *OrderService) GetOrderStatus(ctx context.Context, orderId string) error {
	path := "/v1/orders/{order_id}/status"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CancelOrder Cancel an order (only allowed before PAID or POLICY_ISSUED)
func (s *OrderService) CancelOrder(ctx context.Context, orderId string, req *models.OrderCancellationRequest) error {
	path := "/v1/orders/{order_id}:cancel"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ConfirmPayment Confirm payment for an order (called by payment gateway callback)
func (s *OrderService) ConfirmPayment(ctx context.Context, orderId string, req *models.PaymentConfirmationRequest) error {
	path := "/v1/orders/{order_id}:confirm-payment"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// InitiatePayment Initiate payment for a pending order
func (s *OrderService) InitiatePayment(ctx context.Context, orderId string, req *models.OrdersInitiatePaymentRequest) error {
	path := "/v1/orders/{order_id}:pay"
	path = strings.ReplaceAll(path, "{order_id}", orderId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

