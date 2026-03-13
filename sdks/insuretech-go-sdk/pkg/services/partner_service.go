package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// PartnerService handles partner-related API calls
type PartnerService struct {
	Client Client
}

// ListPartners List partners
func (s *PartnerService) ListPartners(ctx context.Context) (*models.PartnersListingResponse, error) {
	path := "/v1/partners"
	var result models.PartnersListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePartner Partner Management
func (s *PartnerService) CreatePartner(ctx context.Context, req *models.PartnerCreationRequest) (*models.PartnerCreationResponse, error) {
	path := "/v1/partners"
	var result models.PartnerCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPartnerCommission Partner Commission & Financials
func (s *PartnerService) GetPartnerCommission(ctx context.Context, partnerId string) (*models.PartnerCommissionRetrievalResponse, error) {
	path := "/v1/partners/{partner_id}/commission"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.PartnerCommissionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCommissionStructure Update commission structure
func (s *PartnerService) UpdateCommissionStructure(ctx context.Context, partnerId string, req *models.CommissionStructureUpdateRequest) (*models.CommissionStructureUpdateResponse, error) {
	path := "/v1/partners/{partner_id}/commission"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.CommissionStructureUpdateResponse
	err := s.Client.DoRequest(ctx, "PUT", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyPartner Partner Verification & Onboarding
func (s *PartnerService) VerifyPartner(ctx context.Context, partnerId string, req *models.PartnerVerificationRequest) (*models.PartnerVerificationResponse, error) {
	path := "/v1/partners/{partner_id}:verify"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.PartnerVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RotatePartnerAPIKey Rotate partner a p i key
func (s *PartnerService) RotatePartnerAPIKey(ctx context.Context, partnerId string, req *models.PartnerAPIKeyRotationRequest) (*models.PartnerAPIKeyRotationResponse, error) {
	path := "/v1/partners/{partner_id}/credentials:rotate"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.PartnerAPIKeyRotationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPartnerAPICredentials Partner Integration
func (s *PartnerService) GetPartnerAPICredentials(ctx context.Context, partnerId string) (*models.PartnerAPICredentialsRetrievalResponse, error) {
	path := "/v1/partners/{partner_id}/credentials"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.PartnerAPICredentialsRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePartnerStatus Update partner status
func (s *PartnerService) UpdatePartnerStatus(ctx context.Context, partnerId string, req *models.PartnerStatusUpdateRequest) (*models.PartnerStatusUpdateResponse, error) {
	path := "/v1/partners/{partner_id}:update-status"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.PartnerStatusUpdateResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPartner Get partner
func (s *PartnerService) GetPartner(ctx context.Context, partnerId string) (*models.PartnerRetrievalResponse, error) {
	path := "/v1/partners/{partner_id}"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.PartnerRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdatePartner Update partner
func (s *PartnerService) UpdatePartner(ctx context.Context, partnerId string, req *models.PartnerUpdateRequest) (*models.PartnerUpdateResponse, error) {
	path := "/v1/partners/{partner_id}"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	var result models.PartnerUpdateResponse
	err := s.Client.DoRequest(ctx, "PATCH", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePartner Delete partner
func (s *PartnerService) DeletePartner(ctx context.Context, partnerId string) error {
	path := "/v1/partners/{partner_id}"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

