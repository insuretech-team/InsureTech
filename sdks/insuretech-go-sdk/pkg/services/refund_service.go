package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// RefundService handles refund-related API calls
type RefundService struct {
	Client Client
}

// RequestRefund Request refund
func (s *RefundService) RequestRefund(ctx context.Context, policyId string, req *models.RequestRefundRequest) error {
	path := "/v1/policies/{policy_id}/refund"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// CalculateRefund Calculate refund amount
func (s *RefundService) CalculateRefund(ctx context.Context, policyId string, req *models.RefundCalculationRequest) error {
	path := "/v1/policies/{policy_id}/refunds:calculate"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListRefunds List refunds
func (s *RefundService) ListRefunds(ctx context.Context) error {
	path := "/v1/refunds"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetRefund Get refund
func (s *RefundService) GetRefund(ctx context.Context, refundId string) error {
	path := "/v1/refunds/{refund_id}"
	path = strings.ReplaceAll(path, "{refund_id}", refundId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ApproveRefund Approve refund
func (s *RefundService) ApproveRefund(ctx context.Context, refundId string, req *models.RefundApprovalRequest) error {
	path := "/v1/refunds/{refund_id}:approve"
	path = strings.ReplaceAll(path, "{refund_id}", refundId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ProcessRefund Process refund payment
func (s *RefundService) ProcessRefund(ctx context.Context, refundId string, req *models.RefundProcessingRequest) error {
	path := "/v1/refunds/{refund_id}:process"
	path = strings.ReplaceAll(path, "{refund_id}", refundId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

