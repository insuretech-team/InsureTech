package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	paymententityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/entity/v1"
	paymentservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/payment/services/v1"
	"google.golang.org/grpc"
)

type mockPaymentServiceClient struct {
	webhookFn    func(ctx context.Context, in *paymentservicev1.HandleGatewayWebhookRequest, opts ...grpc.CallOption) (*paymentservicev1.HandleGatewayWebhookResponse, error)
	getPaymentFn func(ctx context.Context, in *paymentservicev1.GetPaymentRequest, opts ...grpc.CallOption) (*paymentservicev1.GetPaymentResponse, error)
}

func (m *mockPaymentServiceClient) InitiatePayment(context.Context, *paymentservicev1.InitiatePaymentRequest, ...grpc.CallOption) (*paymentservicev1.InitiatePaymentResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) VerifyPayment(ctx context.Context, in *paymentservicev1.VerifyPaymentRequest, opts ...grpc.CallOption) (*paymentservicev1.VerifyPaymentResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) GetPayment(ctx context.Context, in *paymentservicev1.GetPaymentRequest, opts ...grpc.CallOption) (*paymentservicev1.GetPaymentResponse, error) {
	if m.getPaymentFn == nil {
		return &paymentservicev1.GetPaymentResponse{}, nil
	}
	return m.getPaymentFn(ctx, in, opts...)
}
func (m *mockPaymentServiceClient) ListPayments(context.Context, *paymentservicev1.ListPaymentsRequest, ...grpc.CallOption) (*paymentservicev1.ListPaymentsResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) HandleGatewayWebhook(ctx context.Context, in *paymentservicev1.HandleGatewayWebhookRequest, opts ...grpc.CallOption) (*paymentservicev1.HandleGatewayWebhookResponse, error) {
	if m.webhookFn == nil {
		return &paymentservicev1.HandleGatewayWebhookResponse{}, nil
	}
	return m.webhookFn(ctx, in, opts...)
}
func (m *mockPaymentServiceClient) GetPaymentByProviderReference(context.Context, *paymentservicev1.GetPaymentByProviderReferenceRequest, ...grpc.CallOption) (*paymentservicev1.GetPaymentByProviderReferenceResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) SubmitManualPaymentProof(context.Context, *paymentservicev1.SubmitManualPaymentProofRequest, ...grpc.CallOption) (*paymentservicev1.SubmitManualPaymentProofResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) ReviewManualPayment(context.Context, *paymentservicev1.ReviewManualPaymentRequest, ...grpc.CallOption) (*paymentservicev1.ReviewManualPaymentResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) GenerateReceipt(context.Context, *paymentservicev1.GenerateReceiptRequest, ...grpc.CallOption) (*paymentservicev1.GenerateReceiptResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) GetPaymentReceipt(context.Context, *paymentservicev1.GetPaymentReceiptRequest, ...grpc.CallOption) (*paymentservicev1.GetPaymentReceiptResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) InitiateRefund(context.Context, *paymentservicev1.InitiateRefundRequest, ...grpc.CallOption) (*paymentservicev1.InitiateRefundResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) GetRefundStatus(context.Context, *paymentservicev1.GetRefundStatusRequest, ...grpc.CallOption) (*paymentservicev1.GetRefundStatusResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) ListPaymentMethods(context.Context, *paymentservicev1.ListPaymentMethodsRequest, ...grpc.CallOption) (*paymentservicev1.ListPaymentMethodsResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) AddPaymentMethod(context.Context, *paymentservicev1.AddPaymentMethodRequest, ...grpc.CallOption) (*paymentservicev1.AddPaymentMethodResponse, error) {
	return nil, nil
}
func (m *mockPaymentServiceClient) ReconcilePayments(context.Context, *paymentservicev1.ReconcilePaymentsRequest, ...grpc.CallOption) (*paymentservicev1.ReconcilePaymentsResponse, error) {
	return nil, nil
}

func TestPaymentCallbackHandler_Webhook_MissingIdentifiers_Returns400(t *testing.T) {
	h := &PaymentCallbackHandler{client: &mockPaymentServiceClient{}}
	req := httptest.NewRequest(http.MethodPost, "/v1/payments/webhook/sslcommerz", strings.NewReader("status=VALID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Webhook(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestPaymentCallbackHandler_Webhook_ForwardsToPaymentService(t *testing.T) {
	called := false
	h := &PaymentCallbackHandler{client: &mockPaymentServiceClient{
		webhookFn: func(_ context.Context, in *paymentservicev1.HandleGatewayWebhookRequest, _ ...grpc.CallOption) (*paymentservicev1.HandleGatewayWebhookResponse, error) {
			called = true
			if in.GetProvider() != "sslcommerz" {
				t.Fatalf("unexpected provider %q", in.GetProvider())
			}
			if !strings.Contains(string(in.GetRawPayload()), "value_a=pay-123") {
				t.Fatalf("unexpected payload %q", string(in.GetRawPayload()))
			}
			if in.GetHeaders()["x-payment-callback-type"] != "webhook" {
				t.Fatalf("unexpected callback type %q", in.GetHeaders()["x-payment-callback-type"])
			}
			return &paymentservicev1.HandleGatewayWebhookResponse{
				PaymentId: "pay-123",
				Status:    "PAYMENT_STATUS_SUCCESS",
				Accepted:  true,
			}, nil
		},
	}}
	req := httptest.NewRequest(http.MethodPost, "/v1/payments/webhook/sslcommerz", strings.NewReader("value_a=pay-123&val_id=val-123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Webhook(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !called {
		t.Fatal("expected HandleGatewayWebhook to be called")
	}
}

func TestPaymentCallbackHandler_Success_RedirectsToCallbackURL(t *testing.T) {
	h := &PaymentCallbackHandler{client: &mockPaymentServiceClient{
		webhookFn: func(_ context.Context, in *paymentservicev1.HandleGatewayWebhookRequest, _ ...grpc.CallOption) (*paymentservicev1.HandleGatewayWebhookResponse, error) {
			return &paymentservicev1.HandleGatewayWebhookResponse{
				PaymentId: "pay-123",
				Status:    "PAYMENT_STATUS_SUCCESS",
				Accepted:  true,
			}, nil
		},
		getPaymentFn: func(_ context.Context, in *paymentservicev1.GetPaymentRequest, _ ...grpc.CallOption) (*paymentservicev1.GetPaymentResponse, error) {
			return &paymentservicev1.GetPaymentResponse{
				Payment: &paymententityv1.Payment{
					PaymentId:       in.GetPaymentId(),
					Amount:          &commonv1.Money{Amount: 100, Currency: "BDT"},
					GatewayResponse: `{"callback_url":"https://customer.example.com/payments/result"}`,
				},
			}, nil
		},
	}}
	req := httptest.NewRequest(http.MethodPost, "/v1/payments/sslcommerz/success", strings.NewReader("value_a=pay-123&val_id=val-123&tran_id=tran-123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Success(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}
	location := rec.Header().Get("Location")
	parsed, err := url.Parse(location)
	if err != nil {
		t.Fatalf("invalid redirect url: %v", err)
	}
	if parsed.Host != "customer.example.com" {
		t.Fatalf("unexpected redirect host %q", parsed.Host)
	}
	if parsed.Query().Get("payment_id") != "pay-123" {
		t.Fatalf("unexpected payment_id in redirect")
	}
}

func TestPaymentCallbackHandler_Cancel_ReturnsJSONWithoutCallbackURL(t *testing.T) {
	h := &PaymentCallbackHandler{client: &mockPaymentServiceClient{
		webhookFn: func(_ context.Context, _ *paymentservicev1.HandleGatewayWebhookRequest, _ ...grpc.CallOption) (*paymentservicev1.HandleGatewayWebhookResponse, error) {
			return &paymentservicev1.HandleGatewayWebhookResponse{
				PaymentId: "pay-123",
				Status:    "PAYMENT_STATUS_CANCELLED",
				Accepted:  true,
			}, nil
		},
		getPaymentFn: func(_ context.Context, in *paymentservicev1.GetPaymentRequest, _ ...grpc.CallOption) (*paymentservicev1.GetPaymentResponse, error) {
			return &paymentservicev1.GetPaymentResponse{
				Payment: &paymententityv1.Payment{PaymentId: in.GetPaymentId()},
			}, nil
		},
	}}
	req := httptest.NewRequest(http.MethodPost, "/v1/payments/sslcommerz/cancel", strings.NewReader("value_a=pay-123&tran_id=tran-123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Cancel(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
