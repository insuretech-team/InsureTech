package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// BeneficiaryService handles beneficiary-related API calls
type BeneficiaryService struct {
	Client Client
}

// ListBeneficiaries List beneficiaries (admin)
func (s *BeneficiaryService) ListBeneficiaries(ctx context.Context) error {
	path := "/v1/beneficiaries"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateBusinessBeneficiary Create business beneficiary
func (s *BeneficiaryService) CreateBusinessBeneficiary(ctx context.Context, req *models.BeneficiaryBusinessBeneficiaryCreationRequest) error {
	path := "/v1/beneficiaries/business"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// CreateIndividualBeneficiary Create individual beneficiary
func (s *BeneficiaryService) CreateIndividualBeneficiary(ctx context.Context, req *models.BeneficiaryIndividualBeneficiaryCreationRequest) error {
	path := "/v1/beneficiaries/individual"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetBeneficiary Get beneficiary details
func (s *BeneficiaryService) GetBeneficiary(ctx context.Context, beneficiaryId string) error {
	path := "/v1/beneficiaries/{beneficiary_id}"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateBeneficiary Update beneficiary
func (s *BeneficiaryService) UpdateBeneficiary(ctx context.Context, beneficiaryId string, req *models.BeneficiaryBeneficiaryUpdateRequest) error {
	path := "/v1/beneficiaries/{beneficiary_id}"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// CompleteKYC Complete KYC
func (s *BeneficiaryService) CompleteKYC(ctx context.Context, beneficiaryId string, req *models.KYCCompletionRequest) error {
	path := "/v1/beneficiaries/{beneficiary_id}/kyc"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdateRiskScore Update risk score
func (s *BeneficiaryService) UpdateRiskScore(ctx context.Context, beneficiaryId string, req *models.RiskScoreUpdateRequest) error {
	path := "/v1/beneficiaries/{beneficiary_id}/risk-score"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

