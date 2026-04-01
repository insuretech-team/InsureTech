package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// PaymentService handles payment-related API calls
type PaymentService struct {
	Client Client
}

// ListPayments List payments
func (s *PaymentService) ListPayments(ctx context.Context) error {
	path := "/v1/payments"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// InitiatePayment Payment Processing
func (s *PaymentService) InitiatePayment(ctx context.Context, req *models.PaymentInitiatePaymentRequest) error {
	path := "/v1/payments"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetPaymentByProviderReference Lookup by provider-specific reference (e
func (s *PaymentService) GetPaymentByProviderReference(ctx context.Context, provider string, providerReference string) error {
	path := "/v1/payments/provider/{provider}/references/{provider_reference}"
	path = strings.ReplaceAll(path, "{provider}", provider)
	path = strings.ReplaceAll(path, "{provider_reference}", providerReference)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// HandleGatewayWebhook Gateway webhook — called by API gateway when SSLCommerz/bKash/Nagad posts callback
func (s *PaymentService) HandleGatewayWebhook(ctx context.Context, provider string, req *models.GatewayWebhookHandlingRequest) error {
	path := "/v1/payments/webhook/{provider}"
	path = strings.ReplaceAll(path, "{provider}", provider)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetPayment Get payment
func (s *PaymentService) GetPayment(ctx context.Context, paymentId string) error {
	path := "/v1/payments/{payment_id}"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetPaymentReceipt Retrieve generated receipt (pre-signed URL or file ID)
func (s *PaymentService) GetPaymentReceipt(ctx context.Context, paymentId string) error {
	path := "/v1/payments/{payment_id}/receipt"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// InitiateRefund Refund Management
func (s *PaymentService) InitiateRefund(ctx context.Context, paymentId string, req *models.InitiateRefundRequest) error {
	path := "/v1/payments/{payment_id}/refunds"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GenerateReceipt Trigger async receipt PDF generation after payment is verified
func (s *PaymentService) GenerateReceipt(ctx context.Context, paymentId string, req *models.ReceiptGenerationRequest) error {
	path := "/v1/payments/{payment_id}:generate-receipt"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ReviewManualPayment Admin/agent reviews and approves or rejects a manual payment proof
func (s *PaymentService) ReviewManualPayment(ctx context.Context, paymentId string, req *models.ManualPaymentReviewRequest) error {
	path := "/v1/payments/{payment_id}:review"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SubmitManualPaymentProof Manual bank transfer: customer submits payment proof (scanned deposit slip / screenshot)
func (s *PaymentService) SubmitManualPaymentProof(ctx context.Context, paymentId string, req *models.ManualPaymentProofSubmissionRequest) error {
	path := "/v1/payments/{payment_id}:submit-proof"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// VerifyPayment Verify payment
func (s *PaymentService) VerifyPayment(ctx context.Context, paymentId string, req *models.PaymentVerificationRequest) error {
	path := "/v1/payments/{payment_id}:verify"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ReconcilePayments Reconciliation
func (s *PaymentService) ReconcilePayments(ctx context.Context, req *models.ReconcilePaymentsRequest) error {
	path := "/v1/payments:reconcile"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetRefundStatus Get refund status
func (s *PaymentService) GetRefundStatus(ctx context.Context, refundId string) error {
	path := "/v1/refunds/{refund_id}/status"
	path = strings.ReplaceAll(path, "{refund_id}", refundId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ListPaymentMethods Payment Methods
func (s *PaymentService) ListPaymentMethods(ctx context.Context, userId string) error {
	path := "/v1/users/{user_id}/payment-methods"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// AddPaymentMethod Add payment method
func (s *PaymentService) AddPaymentMethod(ctx context.Context, userId string, req *models.AddPaymentMethodRequest) error {
	path := "/v1/users/{user_id}/payment-methods"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

