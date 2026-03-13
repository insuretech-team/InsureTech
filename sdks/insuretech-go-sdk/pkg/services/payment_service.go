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

// ReconcilePayments Reconciliation
func (s *PaymentService) ReconcilePayments(ctx context.Context, req *models.ReconcilePaymentsRequest) (*models.ReconcilePaymentsResponse, error) {
	path := "/v1/payments:reconcile"
	var result models.ReconcilePaymentsResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListPaymentMethods Payment Methods
func (s *PaymentService) ListPaymentMethods(ctx context.Context, userId string) (*models.PaymentMethodsListingResponse, error) {
	path := "/v1/users/{user_id}/payment-methods"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.PaymentMethodsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddPaymentMethod Add payment method
func (s *PaymentService) AddPaymentMethod(ctx context.Context, userId string, req *models.AddPaymentMethodRequest) (*models.AddPaymentMethodResponse, error) {
	path := "/v1/users/{user_id}/payment-methods"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	var result models.AddPaymentMethodResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ReviewManualPayment Admin/agent reviews and approves or rejects a manual payment proof
func (s *PaymentService) ReviewManualPayment(ctx context.Context, paymentId string, req *models.ManualPaymentReviewRequest) (*models.ManualPaymentReviewResponse, error) {
	path := "/v1/payments/{payment_id}:review"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	var result models.ManualPaymentReviewResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPaymentByProviderReference Lookup by provider-specific reference (e
func (s *PaymentService) GetPaymentByProviderReference(ctx context.Context, provider string, providerReference string) (*models.PaymentByProviderReferenceRetrievalResponse, error) {
	path := "/v1/payments/provider/{provider}/references/{provider_reference}"
	path = strings.ReplaceAll(path, "{provider}", provider)
	path = strings.ReplaceAll(path, "{provider_reference}", providerReference)
	var result models.PaymentByProviderReferenceRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRefundStatus Get refund status
func (s *PaymentService) GetRefundStatus(ctx context.Context, refundId string) (*models.RefundStatusRetrievalResponse, error) {
	path := "/v1/refunds/{refund_id}"
	path = strings.ReplaceAll(path, "{refund_id}", refundId)
	var result models.RefundStatusRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SubmitManualPaymentProof Manual bank transfer: customer submits payment proof (scanned deposit slip / screenshot)
func (s *PaymentService) SubmitManualPaymentProof(ctx context.Context, paymentId string, req *models.ManualPaymentProofSubmissionRequest) (*models.ManualPaymentProofSubmissionResponse, error) {
	path := "/v1/payments/{payment_id}:submit-proof"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	var result models.ManualPaymentProofSubmissionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GenerateReceipt Trigger async receipt PDF generation after payment is verified
func (s *PaymentService) GenerateReceipt(ctx context.Context, paymentId string, req *models.ReceiptGenerationRequest) (*models.ReceiptGenerationResponse, error) {
	path := "/v1/payments/{payment_id}:generate-receipt"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	var result models.ReceiptGenerationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPaymentReceipt Retrieve generated receipt (pre-signed URL or file ID)
func (s *PaymentService) GetPaymentReceipt(ctx context.Context, paymentId string) (*models.PaymentReceiptRetrievalResponse, error) {
	path := "/v1/payments/{payment_id}/receipt"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	var result models.PaymentReceiptRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyPayment Verify payment
func (s *PaymentService) VerifyPayment(ctx context.Context, paymentId string, req *models.PaymentVerificationRequest) (*models.PaymentVerificationResponse, error) {
	path := "/v1/payments/{payment_id}:verify"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	var result models.PaymentVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// HandleGatewayWebhook Gateway webhook — called by API gateway when SSLCommerz/bKash/Nagad posts callback
func (s *PaymentService) HandleGatewayWebhook(ctx context.Context, provider string, req *models.GatewayWebhookHandlingRequest) (*models.GatewayWebhookHandlingResponse, error) {
	path := "/v1/payments/webhook/{provider}"
	path = strings.ReplaceAll(path, "{provider}", provider)
	var result models.GatewayWebhookHandlingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListPayments List payments
func (s *PaymentService) ListPayments(ctx context.Context) (*models.PaymentsListingResponse, error) {
	path := "/v1/payments"
	var result models.PaymentsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InitiatePayment Payment Processing
func (s *PaymentService) InitiatePayment(ctx context.Context, req *models.PaymentInitiatePaymentRequest) (*models.PaymentInitiatePaymentResponse, error) {
	path := "/v1/payments"
	var result models.PaymentInitiatePaymentResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPayment Get payment
func (s *PaymentService) GetPayment(ctx context.Context, paymentId string) (*models.PaymentRetrievalResponse, error) {
	path := "/v1/payments/{payment_id}"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	var result models.PaymentRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InitiateRefund Refund Management
func (s *PaymentService) InitiateRefund(ctx context.Context, paymentId string, req *models.InitiateRefundRequest) (*models.InitiateRefundResponse, error) {
	path := "/v1/payments/{payment_id}/refunds"
	path = strings.ReplaceAll(path, "{payment_id}", paymentId)
	var result models.InitiateRefundResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

