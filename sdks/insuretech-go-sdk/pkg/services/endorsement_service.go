package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// EndorsementService handles endorsement-related API calls
type EndorsementService struct {
	Client Client
}

// GetEndorsement Get endorsement
func (s *EndorsementService) GetEndorsement(ctx context.Context, endorsementId string) error {
	path := "/v1/endorsements/{endorsement_id}"
	path = strings.ReplaceAll(path, "{endorsement_id}", endorsementId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ApproveEndorsement Approve endorsement
func (s *EndorsementService) ApproveEndorsement(ctx context.Context, endorsementId string, req *models.EndorsementApprovalRequest) error {
	path := "/v1/endorsements/{endorsement_id}:approve"
	path = strings.ReplaceAll(path, "{endorsement_id}", endorsementId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RejectEndorsement Reject endorsement
func (s *EndorsementService) RejectEndorsement(ctx context.Context, endorsementId string, req *models.EndorsementRejectionRequest) error {
	path := "/v1/endorsements/{endorsement_id}:reject"
	path = strings.ReplaceAll(path, "{endorsement_id}", endorsementId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListEndorsements List endorsements for policy
func (s *EndorsementService) ListEndorsements(ctx context.Context, policyId string) error {
	path := "/v1/policies/{policy_id}/endorsements"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// RequestEndorsement Request endorsement
func (s *EndorsementService) RequestEndorsement(ctx context.Context, policyId string, req *models.RequestEndorsementRequest) error {
	path := "/v1/policies/{policy_id}/endorsements"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

