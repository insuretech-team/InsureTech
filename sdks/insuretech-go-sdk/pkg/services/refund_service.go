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

// ApproveRefund Approve refund
func (s *RefundService) ApproveRefund(ctx context.Context, refundId string, req *models.RefundApprovalRequest) (*models.RefundApprovalResponse, error) {
	path := "/v1/refunds/{refund_id}"
	path = strings.ReplaceAll(path, "{refund_id}", refundId)
	var result models.RefundApprovalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListRefunds List refunds
func (s *RefundService) ListRefunds(ctx context.Context) (*models.RefundsListingResponse, error) {
	path := "/v1/refunds"
	var result models.RefundsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CalculateRefund Calculate refund amount
func (s *RefundService) CalculateRefund(ctx context.Context, policyId string, req *models.RefundCalculationRequest) (*models.RefundCalculationResponse, error) {
	path := "/v1/policies/{policy_id}/refunds"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.RefundCalculationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RequestRefund Request refund
func (s *RefundService) RequestRefund(ctx context.Context, policyId string, req *models.RequestRefundRequest) (*models.RequestRefundResponse, error) {
	path := "/v1/policies/{policy_id}/refund"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.RequestRefundResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

