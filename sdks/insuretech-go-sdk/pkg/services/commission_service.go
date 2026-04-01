package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// CommissionService handles commission-related API calls
type CommissionService struct {
	Client Client
}

// CreatePayout Create payout batch
func (s *CommissionService) CreatePayout(ctx context.Context, req *models.PayoutCreationRequest) error {
	path := "/v1/commission-payouts"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ProcessPayout Process payout
func (s *CommissionService) ProcessPayout(ctx context.Context, payoutId string, req *models.PayoutProcessingRequest) error {
	path := "/v1/commission-payouts/{payout_id}:process"
	path = strings.ReplaceAll(path, "{payout_id}", payoutId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListCommissions List commissions for recipient
func (s *CommissionService) ListCommissions(ctx context.Context) error {
	path := "/v1/commissions"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetCommission Get commission details
func (s *CommissionService) GetCommission(ctx context.Context, commissionId string) error {
	path := "/v1/commissions/{commission_id}"
	path = strings.ReplaceAll(path, "{commission_id}", commissionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CalculateCommission Calculate and record commission for policy
func (s *CommissionService) CalculateCommission(ctx context.Context, req *models.CommissionCalculationRequest) error {
	path := "/v1/commissions:calculate"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetRevenueShareReport Get revenue share report
func (s *CommissionService) GetRevenueShareReport(ctx context.Context, insurerId string) error {
	path := "/v1/insurers/{insurer_id}/revenue-share"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetCommissionStatement Get commission statement
func (s *CommissionService) GetCommissionStatement(ctx context.Context, recipientId string) error {
	path := "/v1/recipients/{recipient_id}/commission-statement"
	path = strings.ReplaceAll(path, "{recipient_id}", recipientId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

