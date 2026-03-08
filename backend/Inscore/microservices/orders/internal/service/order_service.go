// Package service implements the order lifecycle business logic.
//
// Order lifecycle:
//
//	PENDING → PAYMENT_INITIATED → PAID → POLICY_ISSUED
//	                           ↘ FAILED
//	PENDING / PAYMENT_INITIATED → CANCELLED
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/middleware"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	ordereventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/events/v1"
	orderservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/services/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// OrderServiceImpl implements domain.OrderService.
type OrderServiceImpl struct {
	repo          domain.OrderRepository
	publisher     *events.Publisher
	paymentClient paymentservicev1.PaymentServiceClient // nil = stubbed (dev mode)
}

// Compile-time interface check.
var _ domain.OrderService = (*OrderServiceImpl)(nil)

func NewOrderService(
	repo domain.OrderRepository,
	publisher *events.Publisher,
	paymentClient paymentservicev1.PaymentServiceClient,
) *OrderServiceImpl {
	return &OrderServiceImpl{
		repo:          repo,
		publisher:     publisher,
		paymentClient: paymentClient,
	}
}

// ─── canCancel ────────────────────────────────────────────────────────────────

func canCancel(status ordersv1.OrderStatus) bool {
	return status == ordersv1.OrderStatus_ORDER_STATUS_PENDING ||
		status == ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED
}

// ─── CreateOrder ─────────────────────────────────────────────────────────────

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *orderservicev1.CreateOrderRequest) (*orderservicev1.CreateOrderResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetQuotationId()) == "" {
		return nil, fmt.Errorf("%w: quotation_id is required", ErrInvalidArgument)
	}

	// Extract full caller identity from gateway-set metadata headers.
	rctx := middleware.ExtractRequestContext(ctx)

	// Resolve tenant — prefer request field, fall back to metadata, then env default.
	tenantID := req.GetTenantId()
	if tenantID == "" {
		tenantID = rctx.TenantID
	}
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	// Resolve customer — prefer request field, fall back to authenticated user.
	customerID := strings.TrimSpace(req.GetCustomerId())
	if customerID == "" {
		customerID = rctx.UserID
	}

	// Resolve organisation (B2B) — prefer request field, fall back to gateway header.
	orgID := req.GetOrganisationId()
	if orgID == "" {
		orgID = rctx.OrganisationID
	}

	// Idempotency: return existing order if same key already created one.
	if key := strings.TrimSpace(req.GetIdempotencyKey()); key != "" {
		existing, err := s.repo.GetOrderByIdempotencyKey(ctx, key)
		if err == nil && existing != nil {
			appLogger.Infof("CreateOrder: idempotency hit — returning existing order %s", existing.OrderId)
			return &orderservicev1.CreateOrderResponse{
				Order:   &orderservicev1.OrderView{Order: existing},
				Message: "Order already exists (idempotency replay)",
			}, nil
		}
	}

	// Resolve total payable — prefer explicit request field (caller sourced from quotation).
	totalPayable := req.GetTotalPayable()
	if totalPayable == nil || totalPayable.Amount <= 0 {
		// Fallback: use 1 paisa as sentinel; caller should always supply real amount from quotation.
		// This is intentionally minimal — production callers (PoliSync) must supply total_payable.
		appLogger.Warnf("CreateOrder: total_payable not provided or zero — using 1 paisa sentinel. " +
			"Caller should resolve from quotation before creating order.")
		totalPayable = &commonv1.Money{Amount: 1, Currency: "BDT"}
	}

	// Resolve payment due date — default to 30 minutes from now if not supplied.
	var paymentDueAt *time.Time
	if req.GetPaymentDueAt() != nil {
		t := req.GetPaymentDueAt().AsTime()
		paymentDueAt = &t
	} else {
		t := time.Now().UTC().Add(30 * time.Minute)
		paymentDueAt = &t
	}

	// Coverage timestamps
	var coverageStartAt, coverageEndAt *time.Time
	if req.GetCoverageStartAt() != nil {
		t := req.GetCoverageStartAt().AsTime()
		coverageStartAt = &t
	}
	if req.GetCoverageEndAt() != nil {
		t := req.GetCoverageEndAt().AsTime()
		coverageEndAt = &t
	}

	corrID := uuid.NewString()

	input := domain.OrderCreateInput{
		TenantID:          tenantID,
		QuotationID:       req.GetQuotationId(),
		CustomerID:        customerID,
		ProductID:         req.GetProductId(),
		PlanID:            req.GetPlanId(),
		Currency:          totalPayable.GetCurrency(),
		TotalPayable:      totalPayable,
		OrganisationID:    orgID,
		IdempotencyKey:    req.GetIdempotencyKey(),
		CorrelationID:     corrID,
		PaymentStatus:     ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_UNPAID,
		BillingStatus:     ordersv1.OrderBillingStatus_ORDER_BILLING_STATUS_NOT_INVOICED,
		FulfillmentStatus: ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_NOT_STARTED,
		PaymentDueAt:      paymentDueAt,
		CoverageStartAt:   coverageStartAt,
		CoverageEndAt:     coverageEndAt,
		ActorUserID:       rctx.ActorUserID(),
		Portal:            rctx.Portal,
		PurchaseOrderID:   req.GetPurchaseOrderId(),
	}

	order, err := s.repo.CreateOrder(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	// Publish domain event — non-blocking; failure is logged, not fatal.
	occurredAt := timestamppb.Now()
	if s.publisher != nil {
		s.publisher.Publish(ctx, events.TopicOrderCreated, order.OrderId, &ordereventsv1.OrderCreatedEvent{
			EventId:         uuid.NewString(),
			OrderId:         order.OrderId,
			OrderNumber:     order.OrderNumber,
			QuotationId:     order.QuotationId,
			CustomerId:      order.CustomerId,
			ProductId:       order.ProductId,
			PlanId:          order.PlanId,
			TotalPayable:    order.TotalPayable,
			Timestamp:       occurredAt,
			CorrelationId:   corrID,
			TenantId:        tenantID,
			OrganisationId:  orgID,
			Portal:          rctx.Portal,
			ActorUserId:     rctx.ActorUserID(),
			SessionId:       rctx.SessionID,
			SessionType:     rctx.SessionType,
			IdempotencyKey:  req.GetIdempotencyKey(),
			PurchaseOrderId: req.GetPurchaseOrderId(),
			OccurredAt:      occurredAt,
		})
	}

	return &orderservicev1.CreateOrderResponse{
		Order:   &orderservicev1.OrderView{Order: order},
		Message: "Order created successfully",
	}, nil
}

// ─── GetOrder ────────────────────────────────────────────────────────────────

func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *orderservicev1.GetOrderRequest) (*orderservicev1.GetOrderResponse, error) {
	if req == nil || strings.TrimSpace(req.GetOrderId()) == "" {
		return nil, fmt.Errorf("%w: order_id is required", ErrInvalidArgument)
	}

	order, err := s.repo.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: order %s not found", ErrNotFound, req.GetOrderId())
		}
		return nil, fmt.Errorf("GetOrder: %w", err)
	}

	return &orderservicev1.GetOrderResponse{
		Order: &orderservicev1.OrderView{Order: order},
	}, nil
}

// ─── ListOrders ──────────────────────────────────────────────────────────────

func (s *OrderServiceImpl) ListOrders(ctx context.Context, req *orderservicev1.ListOrdersRequest) (*orderservicev1.ListOrdersResponse, error) {
	if req == nil {
		req = &orderservicev1.ListOrdersRequest{}
	}

	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := 0
	if req.GetPageToken() != "" {
		fmt.Sscanf(req.GetPageToken(), "%d", &offset)
	}

	orders, total, err := s.repo.ListOrders(ctx, pageSize, offset, req.GetCustomerId(), req.GetStatus())
	if err != nil {
		return nil, fmt.Errorf("ListOrders: %w", err)
	}

	views := make([]*orderservicev1.OrderView, 0, len(orders))
	for _, o := range orders {
		views = append(views, &orderservicev1.OrderView{Order: o})
	}

	nextPageToken := ""
	nextOffset := offset + len(orders)
	if int64(nextOffset) < total {
		nextPageToken = fmt.Sprintf("%d", nextOffset)
	}

	return &orderservicev1.ListOrdersResponse{
		Orders:        views,
		NextPageToken: nextPageToken,
		TotalCount:    int32(total),
	}, nil
}

// ─── InitiatePayment ─────────────────────────────────────────────────────────

func (s *OrderServiceImpl) InitiatePayment(ctx context.Context, req *orderservicev1.InitiatePaymentRequest) (*orderservicev1.InitiatePaymentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetOrderId()) == "" {
		return nil, fmt.Errorf("%w: order_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetPaymentMethod()) == "" {
		return nil, fmt.Errorf("%w: payment_method is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetIdempotencyKey()) == "" {
		return nil, fmt.Errorf("%w: idempotency_key is required", ErrInvalidArgument)
	}

	rctx := middleware.ExtractRequestContext(ctx)

	order, err := s.repo.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: order %s not found", ErrNotFound, req.GetOrderId())
		}
		return nil, fmt.Errorf("InitiatePayment fetch order: %w", err)
	}

	// Validate state — only PENDING orders can initiate payment.
	if order.Status != ordersv1.OrderStatus_ORDER_STATUS_PENDING {
		return nil, fmt.Errorf("%w: order is in status %s, expected PENDING",
			ErrInvalidTransition, order.Status.String())
	}

	var (
		paymentID  string
		gatewayRef string
		paymentURL string
		expiresAt  = timestamppb.New(time.Now().Add(30 * time.Minute))
	)

	if s.paymentClient != nil {
		// Real path: call payment-service via gRPC.
		// Use typed fields for core linkage (no more untyped metadata for order_id/tenant_id/etc).
		payResp, err := s.paymentClient.InitiatePayment(ctx, &paymentservicev1.InitiatePaymentRequest{
			UserId:         rctx.UserID,
			PolicyId:       order.PolicyId, // may be empty at this stage for pre-issuance orders
			Amount:         order.TotalPayable,
			Currency:       order.Currency,
			PaymentMethod:  req.GetPaymentMethod(),
			CallbackUrl:    req.GetCallbackUrl(),
			IdempotencyKey: req.GetIdempotencyKey(),
			// Typed linkage fields — replaces metadata map for core order/tenant context
			OrderId:         order.OrderId,
			InvoiceId:       order.InvoiceId,
			TenantId:        rctx.TenantID,
			CustomerId:      order.CustomerId,
			OrganisationId:  rctx.OrganisationID,
			PurchaseOrderId: order.PurchaseOrderId,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: payment-service.InitiatePayment: %v", ErrPaymentFailed, err)
		}
		paymentID = payResp.GetPaymentId()
		// TransactionId doubles as gateway ref until payment proto has GatewayRef field.
		gatewayRef = payResp.GetTransactionId()
		paymentURL = payResp.GetPaymentUrl()
		if payResp.GetExpiresAt() != nil {
			expiresAt = payResp.GetExpiresAt()
		}
	} else {
		return nil, fmt.Errorf("%w: payment-service client not configured", ErrPaymentFailed)
	}

	if err := s.repo.SetPaymentInfo(ctx, req.GetOrderId(), paymentID, gatewayRef); err != nil {
		return nil, fmt.Errorf("InitiatePayment SetPaymentInfo: %w", err)
	}
	if err := s.repo.SetPaymentStatus(ctx, req.GetOrderId(), ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAYMENT_IN_PROGRESS); err != nil {
		appLogger.Warnf("InitiatePayment: SetPaymentStatus failed: %v", err)
	}

	occurredAt := timestamppb.Now()
	if s.publisher != nil {
		s.publisher.Publish(ctx, events.TopicOrderPaymentInitiated, order.OrderId, &ordereventsv1.OrderPaymentInitiatedEvent{
			EventId:           uuid.NewString(),
			OrderId:           order.OrderId,
			PaymentId:         paymentID,
			PaymentGatewayRef: gatewayRef,
			Amount:            order.TotalPayable,
			PaymentMethod:     req.GetPaymentMethod(),
			Timestamp:         occurredAt,
			CorrelationId:     req.GetIdempotencyKey(),
			TenantId:          rctx.TenantID,
			OrganisationId:    rctx.OrganisationID,
			Portal:            rctx.Portal,
			ActorUserId:       rctx.ActorUserID(),
			SessionId:         rctx.SessionID,
			IdempotencyKey:    req.GetIdempotencyKey(),
			OccurredAt:        occurredAt,
		})
	}

	return &orderservicev1.InitiatePaymentResponse{
		OrderId:           req.GetOrderId(),
		PaymentId:         paymentID,
		PaymentUrl:        paymentURL,
		PaymentGatewayRef: gatewayRef,
		ExpiresAt:         expiresAt,
	}, nil
}

// ─── ConfirmPayment ──────────────────────────────────────────────────────────

func (s *OrderServiceImpl) ConfirmPayment(ctx context.Context, req *orderservicev1.ConfirmPaymentRequest) (*orderservicev1.ConfirmPaymentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetOrderId()) == "" {
		return nil, fmt.Errorf("%w: order_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetPaymentId()) == "" {
		return nil, fmt.Errorf("%w: payment_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetTransactionId()) == "" {
		return nil, fmt.Errorf("%w: transaction_id is required", ErrInvalidArgument)
	}

	rctx := middleware.ExtractRequestContext(ctx)

	order, err := s.repo.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: order %s not found", ErrNotFound, req.GetOrderId())
		}
		return nil, fmt.Errorf("ConfirmPayment fetch order: %w", err)
	}

	// Validate state — only PAYMENT_INITIATED orders can be confirmed.
	if order.Status != ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED {
		return nil, fmt.Errorf("%w: order is in status %s, expected PAYMENT_INITIATED",
			ErrInvalidTransition, order.Status.String())
	}

	// Validate payment_id matches what we stored.
	if order.PaymentId != req.GetPaymentId() {
		return nil, fmt.Errorf("%w: payment_id mismatch", ErrInvalidArgument)
	}

	if err := s.repo.UpdateOrderStatus(ctx, req.GetOrderId(), ordersv1.OrderStatus_ORDER_STATUS_PAID); err != nil {
		return nil, fmt.Errorf("ConfirmPayment update status: %w", err)
	}
	if err := s.repo.SetPaymentStatus(ctx, req.GetOrderId(), ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID); err != nil {
		appLogger.Warnf("ConfirmPayment: SetPaymentStatus failed: %v", err)
	}
	if err := s.repo.SetFulfillmentStatus(ctx, req.GetOrderId(), ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLMENT_IN_PROGRESS); err != nil {
		appLogger.Warnf("ConfirmPayment: SetFulfillmentStatus failed: %v", err)
	}

	occurredAt := timestamppb.Now()
	if s.publisher != nil {
		s.publisher.Publish(ctx, events.TopicOrderPaymentConfirmed, order.OrderId, &ordereventsv1.OrderPaymentConfirmedEvent{
			EventId:        uuid.NewString(),
			OrderId:        order.OrderId,
			PaymentId:      req.GetPaymentId(),
			QuotationId:    order.QuotationId,
			CustomerId:     order.CustomerId,
			ProductId:      order.ProductId,
			PlanId:         order.PlanId,
			TotalPayable:   order.TotalPayable,
			Timestamp:      occurredAt,
			CorrelationId:  req.GetTransactionId(),
			TenantId:       rctx.TenantID,
			OrganisationId: rctx.OrganisationID,
			Portal:         rctx.Portal,
			ActorUserId:    rctx.ActorUserID(),
			SessionId:      rctx.SessionID,
			CausationId:    req.GetPaymentId(),
			OccurredAt:     occurredAt,
		})
	}

	return &orderservicev1.ConfirmPaymentResponse{
		OrderId: req.GetOrderId(),
		Status:  ordersv1.OrderStatus_ORDER_STATUS_PAID,
		Message: "Payment confirmed. Policy issuance initiated.",
	}, nil
}

// ─── CancelOrder ─────────────────────────────────────────────────────────────

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, req *orderservicev1.CancelOrderRequest) (*orderservicev1.CancelOrderResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetOrderId()) == "" {
		return nil, fmt.Errorf("%w: order_id is required", ErrInvalidArgument)
	}
	if strings.TrimSpace(req.GetReason()) == "" {
		return nil, fmt.Errorf("%w: reason is required", ErrInvalidArgument)
	}

	rctx := middleware.ExtractRequestContext(ctx)

	order, err := s.repo.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: order %s not found", ErrNotFound, req.GetOrderId())
		}
		return nil, fmt.Errorf("CancelOrder fetch order: %w", err)
	}

	if !canCancel(order.Status) {
		return nil, fmt.Errorf("%w: order in status %s cannot be cancelled",
			ErrInvalidTransition, order.Status.String())
	}

	if err := s.repo.SetCancellationReason(ctx, req.GetOrderId(), req.GetReason()); err != nil {
		return nil, fmt.Errorf("CancelOrder update: %w", err)
	}

	occurredAt := timestamppb.Now()
	if s.publisher != nil {
		s.publisher.Publish(ctx, events.TopicOrderCancelled, order.OrderId, &ordereventsv1.OrderCancelledEvent{
			EventId:        uuid.NewString(),
			OrderId:        order.OrderId,
			CustomerId:     order.CustomerId,
			Reason:         req.GetReason(),
			Timestamp:      occurredAt,
			CorrelationId:  uuid.NewString(),
			TenantId:       rctx.TenantID,
			OrganisationId: rctx.OrganisationID,
			Portal:         rctx.Portal,
			ActorUserId:    rctx.ActorUserID(),
			SessionId:      rctx.SessionID,
			OccurredAt:     occurredAt,
		})
	}

	return &orderservicev1.CancelOrderResponse{
		OrderId: req.GetOrderId(),
		Status:  ordersv1.OrderStatus_ORDER_STATUS_CANCELLED,
		Message: "Order cancelled successfully",
	}, nil
}

// ─── GetOrderStatus ──────────────────────────────────────────────────────────

func (s *OrderServiceImpl) GetOrderStatus(ctx context.Context, req *orderservicev1.GetOrderStatusRequest) (*orderservicev1.GetOrderStatusResponse, error) {
	if req == nil || strings.TrimSpace(req.GetOrderId()) == "" {
		return nil, fmt.Errorf("%w: order_id is required", ErrInvalidArgument)
	}

	order, err := s.repo.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: order %s not found", ErrNotFound, req.GetOrderId())
		}
		return nil, fmt.Errorf("GetOrderStatus: %w", err)
	}

	return &orderservicev1.GetOrderStatusResponse{
		OrderId:   order.OrderId,
		Status:    order.Status,
		PaymentId: order.PaymentId,
		PolicyId:  order.PolicyId,
	}, nil
}

// CreateOrderForB2BPurchaseOrder creates an Order record automatically when a B2B
// PurchaseOrder transitions to APPROVED. This is called by the Kafka consumer so
// the payment flow can begin immediately without manual operator intervention.
//
// The created order carries the purchase_order_id in its reference_id field so
// downstream services (payment, docgen, notification) can correlate the two.
func (s *OrderServiceImpl) CreateOrderForB2BPurchaseOrder(
	ctx context.Context,
	purchaseOrderID string,
	organisationID string,
	tenantID string,
	totalAmount *commonv1.Money,
) error {
	if strings.TrimSpace(purchaseOrderID) == "" {
		return fmt.Errorf("%w: purchase_order_id is required", ErrInvalidArgument)
	}

	req := &orderservicev1.CreateOrderRequest{
		PurchaseOrderId: purchaseOrderID,
		OrganisationId:  organisationID,
		TenantId:        tenantID,
		TotalPayable:    totalAmount,
		PaymentMethod:   "BANK_TRANSFER",      // B2B always starts with bank transfer
		IdempotencyKey:  "b2b-po-" + purchaseOrderID,
	}

	_, err := s.CreateOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateOrderForB2BPurchaseOrder po=%s: %w", purchaseOrderID, err)
	}
	return nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
