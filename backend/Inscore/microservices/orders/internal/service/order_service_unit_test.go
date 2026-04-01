package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/repository"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	orderservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/services/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeOrderRepo struct {
	createInput             domain.OrderCreateInput
	createOrder             *ordersv1.Order
	createErr               error
	getOrder                *ordersv1.Order
	getOrderErr             error
	idempotencyOrder        *ordersv1.Order
	idempotencyErr          error
	listOrders              []*ordersv1.Order
	listTotal               int64
	listErr                 error
	listPageSize            int
	listOffset              int
	listCustomerID          string
	listStatus              ordersv1.OrderStatus
	updateStatusOrderID     string
	updateStatusValue       ordersv1.OrderStatus
	updateStatusErr         error
	setPaymentInfoOrderID   string
	setPaymentInfoPaymentID string
	setPaymentInfoGateway   string
	setPaymentInfoErr       error
	cancelOrderID           string
	cancelReason            string
	cancelErr               error
	failureOrderID          string
	failureReason           string
	failureErr              error
	fulfillmentOrderID      string
	fulfillmentStatus       ordersv1.OrderFulfillmentStatus
	fulfillmentErr          error
	paymentStatusOrderID    string
	paymentStatus           ordersv1.OrderPaymentStatus
	paymentStatusErr        error
	bootstrapPurchaseOrder  *repository.PurchaseOrderBootstrap
	bootstrapErr            error
	ensuredQuotationID      string
	ensureQuotationErr      error
}

func (f *fakeOrderRepo) CreateOrder(_ context.Context, input domain.OrderCreateInput) (*ordersv1.Order, error) {
	f.createInput = input
	if f.createErr != nil {
		return nil, f.createErr
	}
	if f.createOrder != nil {
		return f.createOrder, nil
	}
	return &ordersv1.Order{
		OrderId:         uuid.NewString(),
		OrderNumber:     "ORD-TEST-001",
		QuotationId:     input.QuotationID,
		CustomerId:      input.CustomerID,
		TenantId:        input.TenantID,
		OrganisationId:  input.OrganisationID,
		ProductId:       input.ProductID,
		PlanId:          input.PlanID,
		TotalPayable:    input.TotalPayable,
		Currency:        input.Currency,
		Status:          ordersv1.OrderStatus_ORDER_STATUS_PENDING,
		PurchaseOrderId: input.PurchaseOrderID,
	}, nil
}

func (f *fakeOrderRepo) GetOrder(_ context.Context, _ string) (*ordersv1.Order, error) {
	if f.getOrderErr != nil {
		return nil, f.getOrderErr
	}
	return f.getOrder, nil
}

func (f *fakeOrderRepo) GetOrderByIdempotencyKey(_ context.Context, _ string) (*ordersv1.Order, error) {
	if f.idempotencyErr != nil {
		return nil, f.idempotencyErr
	}
	return f.idempotencyOrder, nil
}

func (f *fakeOrderRepo) ListOrders(_ context.Context, pageSize, offset int, customerID string, status ordersv1.OrderStatus) ([]*ordersv1.Order, int64, error) {
	f.listPageSize = pageSize
	f.listOffset = offset
	f.listCustomerID = customerID
	f.listStatus = status
	return f.listOrders, f.listTotal, f.listErr
}

func (f *fakeOrderRepo) UpdateOrderStatus(_ context.Context, orderID string, status ordersv1.OrderStatus) error {
	f.updateStatusOrderID = orderID
	f.updateStatusValue = status
	return f.updateStatusErr
}

func (f *fakeOrderRepo) SetPaymentInfo(_ context.Context, orderID, paymentID, gatewayRef string) error {
	f.setPaymentInfoOrderID = orderID
	f.setPaymentInfoPaymentID = paymentID
	f.setPaymentInfoGateway = gatewayRef
	return f.setPaymentInfoErr
}

func (f *fakeOrderRepo) SetPolicyID(context.Context, string, string) error  { return nil }
func (f *fakeOrderRepo) SetInvoiceID(context.Context, string, string) error { return nil }

func (f *fakeOrderRepo) SetCancellationReason(_ context.Context, orderID, reason string) error {
	f.cancelOrderID = orderID
	f.cancelReason = reason
	return f.cancelErr
}

func (f *fakeOrderRepo) SetFailureReason(_ context.Context, orderID, reason string) error {
	f.failureOrderID = orderID
	f.failureReason = reason
	return f.failureErr
}

func (f *fakeOrderRepo) SetFulfillmentStatus(_ context.Context, orderID string, status ordersv1.OrderFulfillmentStatus) error {
	f.fulfillmentOrderID = orderID
	f.fulfillmentStatus = status
	return f.fulfillmentErr
}

func (f *fakeOrderRepo) SetPaymentStatus(_ context.Context, orderID string, status ordersv1.OrderPaymentStatus) error {
	f.paymentStatusOrderID = orderID
	f.paymentStatus = status
	return f.paymentStatusErr
}

func (f *fakeOrderRepo) GetPurchaseOrderBootstrap(_ context.Context, _ string) (*repository.PurchaseOrderBootstrap, error) {
	if f.bootstrapErr != nil {
		return nil, f.bootstrapErr
	}
	return f.bootstrapPurchaseOrder, nil
}

func (f *fakeOrderRepo) EnsureApprovedQuotation(_ context.Context, _, _, _ string) (string, error) {
	if f.ensureQuotationErr != nil {
		return "", f.ensureQuotationErr
	}
	if f.ensuredQuotationID != "" {
		return f.ensuredQuotationID, nil
	}
	return "quote-auto-1", nil
}

type fakePaymentClient struct {
	initiateResp *paymentservicev1.InitiatePaymentResponse
	initiateErr  error
	initiateReq  *paymentservicev1.InitiatePaymentRequest
}

func (f *fakePaymentClient) InitiatePayment(_ context.Context, in *paymentservicev1.InitiatePaymentRequest, _ ...grpc.CallOption) (*paymentservicev1.InitiatePaymentResponse, error) {
	f.initiateReq = in
	if f.initiateErr != nil {
		return nil, f.initiateErr
	}
	return f.initiateResp, nil
}

func (f *fakePaymentClient) VerifyPayment(context.Context, *paymentservicev1.VerifyPaymentRequest, ...grpc.CallOption) (*paymentservicev1.VerifyPaymentResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) GetPayment(context.Context, *paymentservicev1.GetPaymentRequest, ...grpc.CallOption) (*paymentservicev1.GetPaymentResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) ListPayments(context.Context, *paymentservicev1.ListPaymentsRequest, ...grpc.CallOption) (*paymentservicev1.ListPaymentsResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) InitiateRefund(context.Context, *paymentservicev1.InitiateRefundRequest, ...grpc.CallOption) (*paymentservicev1.InitiateRefundResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) GetRefundStatus(context.Context, *paymentservicev1.GetRefundStatusRequest, ...grpc.CallOption) (*paymentservicev1.GetRefundStatusResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) ListPaymentMethods(context.Context, *paymentservicev1.ListPaymentMethodsRequest, ...grpc.CallOption) (*paymentservicev1.ListPaymentMethodsResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) AddPaymentMethod(context.Context, *paymentservicev1.AddPaymentMethodRequest, ...grpc.CallOption) (*paymentservicev1.AddPaymentMethodResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) ReconcilePayments(context.Context, *paymentservicev1.ReconcilePaymentsRequest, ...grpc.CallOption) (*paymentservicev1.ReconcilePaymentsResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) HandleGatewayWebhook(context.Context, *paymentservicev1.HandleGatewayWebhookRequest, ...grpc.CallOption) (*paymentservicev1.HandleGatewayWebhookResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) GetPaymentByProviderReference(context.Context, *paymentservicev1.GetPaymentByProviderReferenceRequest, ...grpc.CallOption) (*paymentservicev1.GetPaymentByProviderReferenceResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) SubmitManualPaymentProof(context.Context, *paymentservicev1.SubmitManualPaymentProofRequest, ...grpc.CallOption) (*paymentservicev1.SubmitManualPaymentProofResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) ReviewManualPayment(context.Context, *paymentservicev1.ReviewManualPaymentRequest, ...grpc.CallOption) (*paymentservicev1.ReviewManualPaymentResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) GenerateReceipt(context.Context, *paymentservicev1.GenerateReceiptRequest, ...grpc.CallOption) (*paymentservicev1.GenerateReceiptResponse, error) {
	return nil, nil
}
func (f *fakePaymentClient) GetPaymentReceipt(context.Context, *paymentservicev1.GetPaymentReceiptRequest, ...grpc.CallOption) (*paymentservicev1.GetPaymentReceiptResponse, error) {
	return nil, nil
}

func TestOrderServiceCreateOrderUsesMetadataDefaults(t *testing.T) {
	repo := &fakeOrderRepo{idempotencyErr: domain.ErrNotFound}
	svc := NewOrderService(repo, nil, nil)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-user-id", "user-1",
		"x-tenant-id", "tenant-md",
		"x-portal", "PORTAL_B2B",
		"x-business-id", "org-md",
		"x-session-id", "sess-1",
		"x-session-type", "JWT",
	))

	resp, err := svc.CreateOrder(ctx, &orderservicev1.CreateOrderRequest{
		QuotationId:     "quote-1",
		ProductId:       "product-1",
		PlanId:          "plan-1",
		IdempotencyKey:  "idem-1",
		PurchaseOrderId: "po-1",
	})
	require.NoError(t, err)
	require.Equal(t, "Order created successfully", resp.Message)
	require.Equal(t, "quote-1", repo.createInput.QuotationID)
	require.Equal(t, "user-1", repo.createInput.CustomerID)
	require.Equal(t, "tenant-md", repo.createInput.TenantID)
	require.Equal(t, "org-md", repo.createInput.OrganisationID)
	require.Equal(t, "b2b", repo.createInput.Portal)
	require.Equal(t, "user-1", repo.createInput.ActorUserID)
	require.NotNil(t, repo.createInput.TotalPayable)
	require.Equal(t, int64(1), repo.createInput.TotalPayable.Amount)
	require.Equal(t, "BDT", repo.createInput.TotalPayable.Currency)
	require.NotNil(t, repo.createInput.PaymentDueAt)
	require.NotEmpty(t, repo.createInput.CorrelationID)
}

func TestOrderServiceCreateOrderIdempotencyHit(t *testing.T) {
	existing := &ordersv1.Order{OrderId: "order-1"}
	repo := &fakeOrderRepo{idempotencyOrder: existing}
	svc := NewOrderService(repo, nil, nil)

	resp, err := svc.CreateOrder(context.Background(), &orderservicev1.CreateOrderRequest{
		QuotationId:    "quote-1",
		IdempotencyKey: "same-key",
	})
	require.NoError(t, err)
	require.Equal(t, "Order already exists (idempotency replay)", resp.Message)
	require.Equal(t, "order-1", resp.Order.Order.OrderId)
}

func TestOrderServiceGetListAndStatus(t *testing.T) {
	orderOne := &ordersv1.Order{OrderId: "order-1", Status: ordersv1.OrderStatus_ORDER_STATUS_PENDING}
	orderTwo := &ordersv1.Order{OrderId: "order-2", Status: ordersv1.OrderStatus_ORDER_STATUS_PAID}
	repo := &fakeOrderRepo{
		getOrder:   orderOne,
		listOrders: []*ordersv1.Order{orderOne, orderTwo},
		listTotal:  3,
	}
	svc := NewOrderService(repo, nil, nil)

	getResp, err := svc.GetOrder(context.Background(), &orderservicev1.GetOrderRequest{OrderId: "order-1"})
	require.NoError(t, err)
	require.Equal(t, "order-1", getResp.Order.Order.OrderId)

	listResp, err := svc.ListOrders(context.Background(), &orderservicev1.ListOrdersRequest{
		PageSize:   500,
		PageToken:  "bad-token",
		CustomerId: "customer-1",
		Status:     ordersv1.OrderStatus_ORDER_STATUS_PAID,
	})
	require.NoError(t, err)
	require.Len(t, listResp.Orders, 2)
	require.Equal(t, int32(3), listResp.TotalCount)
	require.Equal(t, "2", listResp.NextPageToken)
	require.Equal(t, 100, repo.listPageSize)
	require.Equal(t, 0, repo.listOffset)
	require.Equal(t, "customer-1", repo.listCustomerID)
	require.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAID, repo.listStatus)

	statusResp, err := svc.GetOrderStatus(context.Background(), &orderservicev1.GetOrderStatusRequest{OrderId: "order-1"})
	require.NoError(t, err)
	require.Equal(t, "order-1", statusResp.OrderId)

	_, err = svc.GetOrder(context.Background(), &orderservicev1.GetOrderRequest{})
	require.ErrorIs(t, err, ErrInvalidArgument)

	repo.getOrderErr = domain.ErrNotFound
	_, err = svc.GetOrderStatus(context.Background(), &orderservicev1.GetOrderStatusRequest{OrderId: "missing"})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestOrderServiceInitiatePaymentSuccess(t *testing.T) {
	repo := &fakeOrderRepo{
		getOrder: &ordersv1.Order{
			OrderId:         "order-1",
			Status:          ordersv1.OrderStatus_ORDER_STATUS_PENDING,
			TotalPayable:    &commonv1.Money{Amount: 5000, Currency: "BDT"},
			Currency:        "BDT",
			CustomerId:      "customer-1",
			InvoiceId:       "invoice-1",
			PurchaseOrderId: "po-1",
		},
	}
	paymentClient := &fakePaymentClient{
		initiateResp: &paymentservicev1.InitiatePaymentResponse{
			PaymentId:     "payment-1",
			TransactionId: "txn-1",
			PaymentUrl:    "https://pay.example/1",
			ExpiresAt:     timestamppb.New(time.Now().Add(time.Hour)),
		},
	}
	svc := NewOrderService(repo, nil, paymentClient)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-user-id", "user-1",
		"x-tenant-id", "tenant-1",
		"x-business-id", "org-1",
		"x-portal", "PORTAL_B2B",
		"x-session-id", "sess-1",
	))

	resp, err := svc.InitiatePayment(ctx, &orderservicev1.InitiatePaymentRequest{
		OrderId:        "order-1",
		PaymentMethod:  "BKASH",
		CallbackUrl:    "https://callback",
		IdempotencyKey: "idem-pay-1",
	})
	require.NoError(t, err)
	require.Equal(t, "payment-1", resp.PaymentId)
	require.Equal(t, "txn-1", resp.PaymentGatewayRef)
	require.Equal(t, "order-1", repo.setPaymentInfoOrderID)
	require.Equal(t, "payment-1", repo.setPaymentInfoPaymentID)
	require.Equal(t, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAYMENT_IN_PROGRESS, repo.paymentStatus)
	require.NotNil(t, paymentClient.initiateReq)
	require.Equal(t, "order-1", paymentClient.initiateReq.OrderId)
	require.Equal(t, "tenant-1", paymentClient.initiateReq.TenantId)
	require.Equal(t, "org-1", paymentClient.initiateReq.OrganisationId)
}

func TestOrderServiceInitiatePaymentErrors(t *testing.T) {
	repo := &fakeOrderRepo{
		getOrder: &ordersv1.Order{
			OrderId: "order-1",
			Status:  ordersv1.OrderStatus_ORDER_STATUS_PAID,
		},
	}
	svc := NewOrderService(repo, nil, nil)

	_, err := svc.InitiatePayment(context.Background(), nil)
	require.ErrorIs(t, err, ErrInvalidArgument)

	_, err = svc.InitiatePayment(context.Background(), &orderservicev1.InitiatePaymentRequest{
		OrderId:       "order-1",
		PaymentMethod: "BKASH",
	})
	require.ErrorIs(t, err, ErrInvalidArgument)

	_, err = svc.InitiatePayment(context.Background(), &orderservicev1.InitiatePaymentRequest{
		OrderId:        "order-1",
		PaymentMethod:  "BKASH",
		IdempotencyKey: "idem",
	})
	require.ErrorIs(t, err, ErrInvalidTransition)

	repo.getOrder.Status = ordersv1.OrderStatus_ORDER_STATUS_PENDING
	_, err = svc.InitiatePayment(context.Background(), &orderservicev1.InitiatePaymentRequest{
		OrderId:        "order-1",
		PaymentMethod:  "BKASH",
		IdempotencyKey: "idem",
	})
	require.ErrorIs(t, err, ErrPaymentFailed)
}

func TestOrderServiceConfirmPaymentAndCancel(t *testing.T) {
	repo := &fakeOrderRepo{
		getOrder: &ordersv1.Order{
			OrderId:      "order-1",
			Status:       ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED,
			PaymentId:    "payment-1",
			QuotationId:  "quote-1",
			CustomerId:   "customer-1",
			ProductId:    "product-1",
			PlanId:       "plan-1",
			TotalPayable: &commonv1.Money{Amount: 1000, Currency: "BDT"},
		},
	}
	svc := NewOrderService(repo, nil, nil)

	confirmResp, err := svc.ConfirmPayment(context.Background(), &orderservicev1.ConfirmPaymentRequest{
		OrderId:       "order-1",
		PaymentId:     "payment-1",
		TransactionId: "txn-1",
	})
	require.NoError(t, err)
	require.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAID, confirmResp.Status)
	require.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_PAID, repo.updateStatusValue)
	require.Equal(t, ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID, repo.paymentStatus)
	require.Equal(t, ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLMENT_IN_PROGRESS, repo.fulfillmentStatus)

	repo.getOrder.PaymentId = "different"
	_, err = svc.ConfirmPayment(context.Background(), &orderservicev1.ConfirmPaymentRequest{
		OrderId:       "order-1",
		PaymentId:     "payment-1",
		TransactionId: "txn-1",
	})
	require.ErrorIs(t, err, ErrInvalidArgument)

	repo.getOrder = &ordersv1.Order{OrderId: "order-2", Status: ordersv1.OrderStatus_ORDER_STATUS_PENDING}
	cancelResp, err := svc.CancelOrder(context.Background(), &orderservicev1.CancelOrderRequest{
		OrderId: "order-2",
		Reason:  "user requested",
	})
	require.NoError(t, err)
	require.Equal(t, ordersv1.OrderStatus_ORDER_STATUS_CANCELLED, cancelResp.Status)
	require.Equal(t, "order-2", repo.cancelOrderID)
	require.Equal(t, "user requested", repo.cancelReason)

	repo.getOrder.Status = ordersv1.OrderStatus_ORDER_STATUS_POLICY_ISSUED
	_, err = svc.CancelOrder(context.Background(), &orderservicev1.CancelOrderRequest{
		OrderId: "order-2",
		Reason:  "too late",
	})
	require.ErrorIs(t, err, ErrInvalidTransition)
}

func TestOrderServiceCreateOrderForB2BPurchaseOrderAndHelpers(t *testing.T) {
	repo := &fakeOrderRepo{
		idempotencyErr: domain.ErrNotFound,
		bootstrapPurchaseOrder: &repository.PurchaseOrderBootstrap{
			BusinessID:       "org-1",
			ProductID:        "product-1",
			PlanID:           "plan-1",
			RequestedBy:      "user-1",
			EstimatedPremium: &commonv1.Money{Amount: 99, Currency: "BDT"},
		},
		ensuredQuotationID: "quote-auto-1",
	}
	svc := NewOrderService(repo, nil, nil)

	err := svc.CreateOrderForB2BPurchaseOrder(context.Background(), "", "org-1", "tenant-1", &commonv1.Money{Amount: 99, Currency: "BDT"})
	require.ErrorIs(t, err, ErrInvalidArgument)

	err = svc.CreateOrderForB2BPurchaseOrder(context.Background(), "po-1", "org-1", "tenant-1", &commonv1.Money{Amount: 99, Currency: "BDT"})
	require.NoError(t, err)
	require.Equal(t, "quote-auto-1", repo.createInput.QuotationID)
	require.Equal(t, "user-1", repo.createInput.CustomerID)
	require.Equal(t, "product-1", repo.createInput.ProductID)
	require.Equal(t, "plan-1", repo.createInput.PlanID)
	require.Equal(t, "po-1", repo.createInput.PurchaseOrderID)

	require.True(t, canCancel(ordersv1.OrderStatus_ORDER_STATUS_PENDING))
	require.True(t, canCancel(ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED))
	require.False(t, canCancel(ordersv1.OrderStatus_ORDER_STATUS_PAID))
	require.Equal(t, 2, min(2, 5))
	require.Equal(t, 5, min(8, 5))
}

func TestOrderServiceCreateOrderPropagatesRepoError(t *testing.T) {
	repo := &fakeOrderRepo{
		idempotencyErr: domain.ErrNotFound,
		createErr:      errors.New("insert failed"),
	}
	svc := NewOrderService(repo, nil, nil)

	_, err := svc.CreateOrder(context.Background(), &orderservicev1.CreateOrderRequest{QuotationId: "quote-1"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "CreateOrder")
}
