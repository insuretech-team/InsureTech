package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// MfsService handles mfs-related API calls
type MfsService struct {
	Client Client
}

// InitiatePayment Initiate payment
func (s *MfsService) InitiatePayment(ctx context.Context, req *models.MfsInitiatePaymentRequest) error {
	path := "/v1/mfs/payments"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ExecuteRefund Execute refund
func (s *MfsService) ExecuteRefund(ctx context.Context, req *models.RefundExecutionRequest) error {
	path := "/v1/mfs/refunds"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListTransactions List transactions
func (s *MfsService) ListTransactions(ctx context.Context) error {
	path := "/v1/mfs/transactions"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetTransaction Get transaction
func (s *MfsService) GetTransaction(ctx context.Context, mfsTransactionId string) error {
	path := "/v1/mfs/transactions/{mfs_transaction_id}"
	path = strings.ReplaceAll(path, "{mfs_transaction_id}", mfsTransactionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ProcessWebhook Process webhook
func (s *MfsService) ProcessWebhook(ctx context.Context, provider string, req *models.WebhookProcessingRequest) error {
	path := "/v1/mfs/webhooks/{provider}"
	path = strings.ReplaceAll(path, "{provider}", provider)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

