package domain

import (
	"context"
	"errors"
	"time"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	orderservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/services/v1"
)

// ─── Sentinel errors ─────────────────────────────────────────────────────────

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists") // idempotency key collision
)

// ─── INPUT TYPES ─────────────────────────────────────────────────────────────

// OrderCreateInput carries all fields needed to create an order row.
type OrderCreateInput struct {
	OrderID         string
	OrderNumber     string
	TenantID        string
	QuotationID     string
	CustomerID      string
	ProductID       string
	PlanID          string
	TotalPayable    *commonv1.Money
	Currency        string
	Status          ordersv1.OrderStatus
	// Phase 2 extended fields
	InvoiceID       string
	OrganisationID  string
	IdempotencyKey  string
	CorrelationID   string
	PaymentStatus   ordersv1.OrderPaymentStatus
	BillingStatus   ordersv1.OrderBillingStatus
	FulfillmentStatus ordersv1.OrderFulfillmentStatus
	PaymentDueAt    *time.Time
	CoverageStartAt *time.Time
	CoverageEndAt   *time.Time
	ManualReviewRequired bool
	ActorUserID          string
	Portal               string
	PurchaseOrderID      string
}

// OrderUpdateInput carries fields for a generic order update.
type OrderUpdateInput struct {
	OrderID            string
	Status             ordersv1.OrderStatus
	PaymentID          string
	PaymentGatewayRef  string
	PolicyID           string
	CancellationReason string
	FailureReason      string
}

// ─── REPOSITORY INTERFACE ────────────────────────────────────────────────────

type OrderRepository interface {
	// CRUD
	CreateOrder(ctx context.Context, input OrderCreateInput) (*ordersv1.Order, error)
	GetOrder(ctx context.Context, orderID string) (*ordersv1.Order, error)
	GetOrderByIdempotencyKey(ctx context.Context, idempotencyKey string) (*ordersv1.Order, error)
	ListOrders(ctx context.Context, pageSize, offset int, customerID string, status ordersv1.OrderStatus) ([]*ordersv1.Order, int64, error)

	// State transitions
	UpdateOrderStatus(ctx context.Context, orderID string, status ordersv1.OrderStatus) error
	SetPaymentInfo(ctx context.Context, orderID, paymentID, gatewayRef string) error
	SetPolicyID(ctx context.Context, orderID, policyID string) error
	SetInvoiceID(ctx context.Context, orderID, invoiceID string) error
	SetCancellationReason(ctx context.Context, orderID, reason string) error
	SetFailureReason(ctx context.Context, orderID, reason string) error
	SetFulfillmentStatus(ctx context.Context, orderID string, status ordersv1.OrderFulfillmentStatus) error
	SetPaymentStatus(ctx context.Context, orderID string, status ordersv1.OrderPaymentStatus) error
}

// ─── SERVICE INTERFACE ───────────────────────────────────────────────────────

type OrderService interface {
	CreateOrder(ctx context.Context, req *orderservicev1.CreateOrderRequest) (*orderservicev1.CreateOrderResponse, error)
	GetOrder(ctx context.Context, req *orderservicev1.GetOrderRequest) (*orderservicev1.GetOrderResponse, error)
	ListOrders(ctx context.Context, req *orderservicev1.ListOrdersRequest) (*orderservicev1.ListOrdersResponse, error)
	InitiatePayment(ctx context.Context, req *orderservicev1.InitiatePaymentRequest) (*orderservicev1.InitiatePaymentResponse, error)
	ConfirmPayment(ctx context.Context, req *orderservicev1.ConfirmPaymentRequest) (*orderservicev1.ConfirmPaymentResponse, error)
	CancelOrder(ctx context.Context, req *orderservicev1.CancelOrderRequest) (*orderservicev1.CancelOrderResponse, error)
	GetOrderStatus(ctx context.Context, req *orderservicev1.GetOrderStatusRequest) (*orderservicev1.GetOrderStatusResponse, error)
}
