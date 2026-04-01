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
func (s *PartnerService) ListPartners(ctx context.Context) error {
	path := "/v1/partners"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreatePartner Partner Management
func (s *PartnerService) CreatePartner(ctx context.Context, req *models.PartnerCreationRequest) error {
	path := "/v1/partners"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetPartner Get partner
func (s *PartnerService) GetPartner(ctx context.Context, partnerId string) error {
	path := "/v1/partners/{partner_id}"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdatePartner Update partner
func (s *PartnerService) UpdatePartner(ctx context.Context, partnerId string, req *models.PartnerUpdateRequest) error {
	path := "/v1/partners/{partner_id}"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// DeletePartner Delete partner
func (s *PartnerService) DeletePartner(ctx context.Context, partnerId string) error {
	path := "/v1/partners/{partner_id}"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "DELETE", path, nil, nil)
}

// GetPartnerCommission Partner Commission & Financials
func (s *PartnerService) GetPartnerCommission(ctx context.Context, partnerId string) error {
	path := "/v1/partners/{partner_id}/commission"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateCommissionStructure Update commission structure
func (s *PartnerService) UpdateCommissionStructure(ctx context.Context, partnerId string, req *models.CommissionStructureUpdateRequest) error {
	path := "/v1/partners/{partner_id}/commission"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "PUT", path, req, nil)
}

// GetPartnerAPICredentials Partner Integration
func (s *PartnerService) GetPartnerAPICredentials(ctx context.Context, partnerId string) error {
	path := "/v1/partners/{partner_id}/credentials"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RotatePartnerAPIKey Rotate partner API key
func (s *PartnerService) RotatePartnerAPIKey(ctx context.Context, partnerId string, req *models.PartnerAPIKeyRotationRequest) error {
	path := "/v1/partners/{partner_id}/credentials:rotate"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// UpdatePartnerStatus Update partner status
func (s *PartnerService) UpdatePartnerStatus(ctx context.Context, partnerId string, req *models.PartnerStatusUpdateRequest) error {
	path := "/v1/partners/{partner_id}:update-status"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// VerifyPartner Partner Verification & Onboarding
func (s *PartnerService) VerifyPartner(ctx context.Context, partnerId string, req *models.PartnerVerificationRequest) error {
	path := "/v1/partners/{partner_id}:verify"
	path = strings.ReplaceAll(path, "{partner_id}", partnerId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

