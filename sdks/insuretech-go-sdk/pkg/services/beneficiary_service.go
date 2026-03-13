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

// CreateIndividualBeneficiary Create individual beneficiary
func (s *BeneficiaryService) CreateIndividualBeneficiary(ctx context.Context, req *models.BeneficiaryIndividualBeneficiaryCreationRequest) (*models.BeneficiaryIndividualBeneficiaryCreationResponse, error) {
	path := "/v1/beneficiaries/individual"
	var result models.BeneficiaryIndividualBeneficiaryCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListBeneficiaries List beneficiaries (admin)
func (s *BeneficiaryService) ListBeneficiaries(ctx context.Context) (*models.BeneficiaryBeneficiariesListingResponse, error) {
	path := "/v1/beneficiaries"
	var result models.BeneficiaryBeneficiariesListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CompleteKYC Complete KYC
func (s *BeneficiaryService) CompleteKYC(ctx context.Context, beneficiaryId string, req *models.KYCCompletionRequest) (*models.KYCCompletionResponse, error) {
	path := "/v1/beneficiaries/{beneficiary_id}/kyc"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	var result models.KYCCompletionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateRiskScore Update risk score
func (s *BeneficiaryService) UpdateRiskScore(ctx context.Context, beneficiaryId string, req *models.RiskScoreUpdateRequest) (*models.RiskScoreUpdateResponse, error) {
	path := "/v1/beneficiaries/{beneficiary_id}/risk-score"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	var result models.RiskScoreUpdateResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateBusinessBeneficiary Create business beneficiary
func (s *BeneficiaryService) CreateBusinessBeneficiary(ctx context.Context, req *models.BeneficiaryBusinessBeneficiaryCreationRequest) (*models.BeneficiaryBusinessBeneficiaryCreationResponse, error) {
	path := "/v1/beneficiaries/business"
	var result models.BeneficiaryBusinessBeneficiaryCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBeneficiary Get beneficiary details
func (s *BeneficiaryService) GetBeneficiary(ctx context.Context, beneficiaryId string) (*models.BeneficiaryBeneficiaryRetrievalResponse, error) {
	path := "/v1/beneficiaries/{beneficiary_id}"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	var result models.BeneficiaryBeneficiaryRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateBeneficiary Update beneficiary
func (s *BeneficiaryService) UpdateBeneficiary(ctx context.Context, beneficiaryId string, req *models.BeneficiaryBeneficiaryUpdateRequest) (*models.BeneficiaryBeneficiaryUpdateResponse, error) {
	path := "/v1/beneficiaries/{beneficiary_id}"
	path = strings.ReplaceAll(path, "{beneficiary_id}", beneficiaryId)
	var result models.BeneficiaryBeneficiaryUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

