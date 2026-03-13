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
func (s *EndorsementService) GetEndorsement(ctx context.Context, endorsementId string) (*models.EndorsementEndorsementRetrievalResponse, error) {
	path := "/v1/endorsements/{endorsement_id}"
	path = strings.ReplaceAll(path, "{endorsement_id}", endorsementId)
	var result models.EndorsementEndorsementRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ApproveEndorsement Approve endorsement
func (s *EndorsementService) ApproveEndorsement(ctx context.Context, endorsementId string, req *models.EndorsementApprovalRequest) (*models.EndorsementApprovalResponse, error) {
	path := "/v1/endorsements/{endorsement_id}"
	path = strings.ReplaceAll(path, "{endorsement_id}", endorsementId)
	var result models.EndorsementApprovalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RequestEndorsement Request endorsement
func (s *EndorsementService) RequestEndorsement(ctx context.Context, policyId string, req *models.RequestEndorsementRequest) (*models.RequestEndorsementResponse, error) {
	path := "/v1/policies/{policy_id}/endorsements"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.RequestEndorsementResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListEndorsements List endorsements for policy
func (s *EndorsementService) ListEndorsements(ctx context.Context, policyId string) (*models.EndorsementsListingResponse, error) {
	path := "/v1/policies/{policy_id}/endorsements"
	path = strings.ReplaceAll(path, "{policy_id}", policyId)
	var result models.EndorsementsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

