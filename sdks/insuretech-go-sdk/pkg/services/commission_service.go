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

// GetRevenueShareReport Get revenue share report
func (s *CommissionService) GetRevenueShareReport(ctx context.Context, insurerId string) (*models.RevenueShareReportRetrievalResponse, error) {
	path := "/v1/insurers/{insurer_id}/revenue-share"
	path = strings.ReplaceAll(path, "{insurer_id}", insurerId)
	var result models.RevenueShareReportRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCommission Get commission details
func (s *CommissionService) GetCommission(ctx context.Context, commissionId string) (*models.CommissionRetrievalResponse, error) {
	path := "/v1/commissions/{commission_id}"
	path = strings.ReplaceAll(path, "{commission_id}", commissionId)
	var result models.CommissionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCommissionStatement Get commission statement
func (s *CommissionService) GetCommissionStatement(ctx context.Context, recipientId string) (*models.CommissionStatementRetrievalResponse, error) {
	path := "/v1/recipients/{recipient_id}/commission-statement"
	path = strings.ReplaceAll(path, "{recipient_id}", recipientId)
	var result models.CommissionStatementRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ProcessPayout Process payout
func (s *CommissionService) ProcessPayout(ctx context.Context, payoutId string, req *models.PayoutProcessingRequest) (*models.PayoutProcessingResponse, error) {
	path := "/v1/commission-payouts/{payout_id}"
	path = strings.ReplaceAll(path, "{payout_id}", payoutId)
	var result models.PayoutProcessingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePayout Create payout batch
func (s *CommissionService) CreatePayout(ctx context.Context, req *models.PayoutCreationRequest) (*models.PayoutCreationResponse, error) {
	path := "/v1/commission-payouts"
	var result models.PayoutCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CalculateCommission Calculate and record commission for policy
func (s *CommissionService) CalculateCommission(ctx context.Context, req *models.CommissionCalculationRequest) (*models.CommissionCalculationResponse, error) {
	path := "/v1/commissions"
	var result models.CommissionCalculationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListCommissions List commissions for recipient
func (s *CommissionService) ListCommissions(ctx context.Context) (*models.CommissionsListingResponse, error) {
	path := "/v1/commissions"
	var result models.CommissionsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

