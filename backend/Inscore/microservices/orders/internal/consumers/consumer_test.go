package consumers

import (
	"context"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/service"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2beventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/events/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	ordersv1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/entity/v1"
	paymenteventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type fakeConsumerRepo struct {
	createInput            domain.OrderCreateInput
	updateStatusOrderID    string
	updateStatus           ordersv1.OrderStatus
	updateStatusErr        error
	paymentStatusOrderID   string
	paymentStatus          ordersv1.OrderPaymentStatus
	paymentStatusErr       error
	failureOrderID         string
	failureReason          string
	failureErr             error
	policyOrderID          string
	policyID               string
	policyErr              error
	fulfillmentOrderID     string
	fulfillmentStatus      ordersv1.OrderFulfillmentStatus
	fulfillmentErr         error
	bootstrapPurchaseOrder *repository.PurchaseOrderBootstrap
	ensureQuotationID      string
}

func (f *fakeConsumerRepo) CreateOrder(_ context.Context, input domain.OrderCreateInput) (*ordersv1.Order, error) {
	f.createInput = input
	return &ordersv1.Order{
		OrderId:         "order-created",
		OrderNumber:     "ORD-TEST-001",
		QuotationId:     input.QuotationID,
		CustomerId:      input.CustomerID,
		ProductId:       input.ProductID,
		PlanId:          input.PlanID,
		PurchaseOrderId: input.PurchaseOrderID,
		TotalPayable:    input.TotalPayable,
		Currency:        input.Currency,
		Status:          ordersv1.OrderStatus_ORDER_STATUS_PENDING,
	}, nil
}
func (f *fakeConsumerRepo) GetOrder(context.Context, string) (*ordersv1.Order, error) {
	return nil, nil
}
func (f *fakeConsumerRepo) GetOrderByIdempotencyKey(context.Context, string) (*ordersv1.Order, error) {
	return nil, nil
}
func (f *fakeConsumerRepo) ListOrders(context.Context, int, int, string, ordersv1.OrderStatus) ([]*ordersv1.Order, int64, error) {
	return nil, 0, nil
}
func (f *fakeConsumerRepo) UpdateOrderStatus(_ context.Context, orderID string, status ordersv1.OrderStatus) error {
	f.updateStatusOrderID = orderID
	f.updateStatus = status
	return f.updateStatusErr
}
func (f *fakeConsumerRepo) SetPaymentInfo(context.Context, string, string, string) error { return nil }
func (f *fakeConsumerRepo) SetPolicyID(_ context.Context, orderID, policyID string) error {
	f.policyOrderID = orderID
	f.policyID = policyID
	return f.policyErr
}
func (f *fakeConsumerRepo) SetInvoiceID(context.Context, string, string) error          { return nil }
func (f *fakeConsumerRepo) SetCancellationReason(context.Context, string, string) error { return nil }
func (f *fakeConsumerRepo) SetFailureReason(_ context.Context, orderID, reason string) error {
	f.failureOrderID = orderID
	f.failureReason = reason
	return f.failureErr
}
func (f *fakeConsumerRepo) SetFulfillmentStatus(_ context.Context, orderID string, status ordersv1.OrderFulfillmentStatus) error {
	f.fulfillmentOrderID = orderID
	f.fulfillmentStatus = status
	return f.fulfillmentErr
}
func (f *fakeConsumerRepo) SetPaymentStatus(_ context.Context, orderID string, status ordersv1.OrderPaymentStatus) error {
	f.paymentStatusOrderID = orderID
	f.paymentStatus = status
	return f.paymentStatusErr
}

func (f *fakeConsumerRepo) GetPurchaseOrderBootstrap(context.Context, string) (*repository.PurchaseOrderBootstrap, error) {
	return f.bootstrapPurchaseOrder, nil
}

func (f *fakeConsumerRepo) EnsureApprovedQuotation(context.Context, string, string, string) (string, error) {
	if f.ensureQuotationID == "" {
		return "quote-auto-1", nil
	}
	return f.ensureQuotationID, nil
}

func TestEventConsumerPaymentAndPolicyHandlers(t *testing.T) {
	repo := &fakeConsumerRepo{}
	consumer := NewEventConsumer(repo, nil, nil, nil)

	completedPayload, _ := protojson.Marshal(&paymenteventsv1.PaymentCompletedEvent{
		OrderId:   "order-1",
		PaymentId: "payment-1",
		Provider:  "sslcommerz",
	})
	if err := consumer.HandlePaymentCompleted(context.Background(), completedPayload); err != nil {
		t.Fatalf("HandlePaymentCompleted() error = %v", err)
	}
	if repo.updateStatusOrderID != "order-1" || repo.updateStatus != ordersv1.OrderStatus_ORDER_STATUS_PAID {
		t.Fatalf("unexpected payment completed update: %+v", repo)
	}
	if repo.paymentStatus != ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID {
		t.Fatalf("unexpected payment status: %+v", repo.paymentStatus)
	}

	failedPayload, _ := protojson.Marshal(&paymenteventsv1.PaymentFailedEvent{
		OrderId:   "order-2",
		ErrorCode: "FAILED",
	})
	if err := consumer.HandlePaymentFailed(context.Background(), failedPayload); err != nil {
		t.Fatalf("HandlePaymentFailed() error = %v", err)
	}
	if repo.failureOrderID != "order-2" || repo.failureReason != "FAILED" {
		t.Fatalf("unexpected failure state: %+v", repo)
	}

	if err := consumer.HandlePolicyIssued(context.Background(), []byte(`{"order_id":"order-3","policy_id":"policy-1"}`)); err != nil {
		t.Fatalf("HandlePolicyIssued() error = %v", err)
	}
	if repo.policyOrderID != "order-3" || repo.policyID != "policy-1" {
		t.Fatalf("unexpected policy link: %+v", repo)
	}
	if repo.fulfillmentStatus != ordersv1.OrderFulfillmentStatus_ORDER_FULFILLMENT_STATUS_FULFILLED {
		t.Fatalf("unexpected fulfillment status: %v", repo.fulfillmentStatus)
	}
}

func TestEventConsumerManualReviewAndVerification(t *testing.T) {
	repo := &fakeConsumerRepo{}
	consumer := NewEventConsumer(repo, nil, nil, nil)

	verifiedPayload, _ := protojson.Marshal(&paymenteventsv1.PaymentVerifiedEvent{
		OrderId:   "order-4",
		PaymentId: "payment-4",
	})
	if err := consumer.HandlePaymentVerified(context.Background(), verifiedPayload); err != nil {
		t.Fatalf("HandlePaymentVerified() error = %v", err)
	}
	if repo.updateStatusOrderID != "order-4" || repo.paymentStatus != ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID {
		t.Fatalf("unexpected verification update: %+v", repo)
	}

	manualReqPayload, _ := protojson.Marshal(&paymenteventsv1.ManualPaymentProofSubmittedEvent{
		OrderId:           "order-5",
		PaymentId:         "payment-5",
		ManualProofFileId: "file-1",
	})
	if err := consumer.HandleManualReviewRequested(context.Background(), manualReqPayload); err != nil {
		t.Fatalf("HandleManualReviewRequested() error = %v", err)
	}
	if repo.updateStatusOrderID != "order-5" || repo.updateStatus != ordersv1.OrderStatus_ORDER_STATUS_PAYMENT_INITIATED {
		t.Fatalf("unexpected manual review pending state: %+v", repo)
	}

	approvedPayload, _ := protojson.Marshal(&paymenteventsv1.ManualPaymentReviewedEvent{
		OrderId:   "order-6",
		PaymentId: "payment-6",
		Approved:  true,
	})
	if err := consumer.HandleManualPaymentReviewed(context.Background(), approvedPayload); err != nil {
		t.Fatalf("HandleManualPaymentReviewed(approve) error = %v", err)
	}
	if repo.updateStatusOrderID != "order-6" || repo.paymentStatus != ordersv1.OrderPaymentStatus_ORDER_PAYMENT_STATUS_PAID {
		t.Fatalf("unexpected approved review state: %+v", repo)
	}

	rejectedPayload, _ := protojson.Marshal(&paymenteventsv1.ManualPaymentReviewedEvent{
		OrderId:         "order-7",
		PaymentId:       "payment-7",
		Approved:        false,
		RejectionReason: "bad proof",
	})
	if err := consumer.HandleManualPaymentReviewed(context.Background(), rejectedPayload); err != nil {
		t.Fatalf("HandleManualPaymentReviewed(reject) error = %v", err)
	}
	if repo.failureOrderID != "order-7" || repo.failureReason != "bad proof" {
		t.Fatalf("unexpected rejected review state: %+v", repo)
	}
}

func TestEventConsumerB2BAndHelpers(t *testing.T) {
	consumer := NewEventConsumer(&fakeConsumerRepo{}, nil, nil, nil)

	ignoredPayload, _ := protojson.Marshal(&b2beventsv1.PurchaseOrderStatusChangedEvent{
		PurchaseOrderId: "po-1",
		OrganisationId:  "org-1",
		NewStatus:       b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED,
	})
	if err := consumer.HandleB2BPurchaseOrderApproved(context.Background(), ignoredPayload); err != nil {
		t.Fatalf("HandleB2BPurchaseOrderApproved(ignore) error = %v", err)
	}

	approvedPayload, _ := protojson.Marshal(&b2beventsv1.PurchaseOrderStatusChangedEvent{
		PurchaseOrderId: "po-2",
		OrganisationId:  "org-1",
		NewStatus:       b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_APPROVED,
	})
	if err := consumer.HandleB2BPurchaseOrderApproved(context.Background(), approvedPayload); err != nil {
		t.Fatalf("HandleB2BPurchaseOrderApproved(nil svc) error = %v", err)
	}

	flat, err := flattenJSONMap([]byte(`{"order_id":"order-1","policy_id":"policy-1","approved":true,"amount":15}`))
	if err != nil {
		t.Fatalf("flattenJSONMap() error = %v", err)
	}
	if flat["approved"] != "true" || flat["amount"] != "15" {
		t.Fatalf("unexpected flattened map: %+v", flat)
	}
}

func TestEventConsumerApprovedPurchaseOrderCreatesOrder(t *testing.T) {
	repo := &fakeConsumerRepo{
		bootstrapPurchaseOrder: &repository.PurchaseOrderBootstrap{
			BusinessID:       "org-1",
			ProductID:        "product-1",
			PlanID:           "plan-1",
			RequestedBy:      "user-1",
			EstimatedPremium: &commonv1.Money{Amount: 4200, Currency: "BDT"},
		},
		ensureQuotationID: "quote-auto-9",
	}
	orderSvc := service.NewOrderService(repo, nil, nil)
	consumer := NewEventConsumer(repo, orderSvc, nil, nil)

	approvedPayload, _ := protojson.Marshal(&b2beventsv1.PurchaseOrderStatusChangedEvent{
		PurchaseOrderId: "po-9",
		OrganisationId:  "org-1",
		NewStatus:       b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_APPROVED,
	})
	if err := consumer.HandleB2BPurchaseOrderApproved(context.Background(), approvedPayload); err != nil {
		t.Fatalf("HandleB2BPurchaseOrderApproved(create) error = %v", err)
	}

	if repo.createInput.QuotationID != "quote-auto-9" || repo.createInput.CustomerID != "user-1" {
		t.Fatalf("unexpected create input: %+v", repo.createInput)
	}
	if repo.createInput.PurchaseOrderID != "po-9" || repo.createInput.ProductID != "product-1" || repo.createInput.PlanID != "plan-1" {
		t.Fatalf("unexpected purchase order mapping: %+v", repo.createInput)
	}
}
