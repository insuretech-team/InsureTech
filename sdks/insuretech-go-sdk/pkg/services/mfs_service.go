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

// GetTransaction Get transaction
func (s *MfsService) GetTransaction(ctx context.Context, mfsTransactionId string) (*models.TransactionRetrievalResponse, error) {
	path := "/v1/mfs/transactions/{mfs_transaction_id}"
	path = strings.ReplaceAll(path, "{mfs_transaction_id}", mfsTransactionId)
	var result models.TransactionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListTransactions List transactions
func (s *MfsService) ListTransactions(ctx context.Context) (*models.TransactionsListingResponse, error) {
	path := "/v1/mfs/transactions"
	var result models.TransactionsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecuteRefund Execute refund
func (s *MfsService) ExecuteRefund(ctx context.Context, req *models.RefundExecutionRequest) (*models.RefundExecutionResponse, error) {
	path := "/v1/mfs/refunds"
	var result models.RefundExecutionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// InitiatePayment Initiate payment
func (s *MfsService) InitiatePayment(ctx context.Context, req *models.MfsInitiatePaymentRequest) (*models.MfsInitiatePaymentResponse, error) {
	path := "/v1/mfs/payments"
	var result models.MfsInitiatePaymentResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ProcessWebhook Process webhook
func (s *MfsService) ProcessWebhook(ctx context.Context, provider string, req *models.WebhookProcessingRequest) (*models.WebhookProcessingResponse, error) {
	path := "/v1/mfs/webhooks/{provider}"
	path = strings.ReplaceAll(path, "{provider}", provider)
	var result models.WebhookProcessingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

