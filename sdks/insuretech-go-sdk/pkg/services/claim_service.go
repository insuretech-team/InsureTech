package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// ClaimService handles claim-related API calls
type ClaimService struct {
	Client Client
}

// DisputeClaim Dispute claim (by customer)
func (s *ClaimService) DisputeClaim(ctx context.Context, claimId string, req *models.DisputeClaimRequest) (*models.DisputeClaimResponse, error) {
	path := "/v1/claims/{claim_id}:dispute"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	var result models.DisputeClaimResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SettleClaim Settle claim
func (s *ClaimService) SettleClaim(ctx context.Context, claimId string, req *models.SettleClaimRequest) (*models.SettleClaimResponse, error) {
	path := "/v1/claims/{claim_id}:settle"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	var result models.SettleClaimResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadDocument Upload claim document
func (s *ClaimService) UploadDocument(ctx context.Context, claimId string, req *models.ClaimsDocumentUploadRequest) (*models.ClaimsDocumentUploadResponse, error) {
	path := "/v1/claims/{claim_id}/documents"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	var result models.ClaimsDocumentUploadResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetClaim Get claim details
func (s *ClaimService) GetClaim(ctx context.Context, claimId string) (*models.ClaimsClaimRetrievalResponse, error) {
	path := "/v1/claims/{claim_id}"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	var result models.ClaimsClaimRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SubmitClaim Submit claim
func (s *ClaimService) SubmitClaim(ctx context.Context, req *models.ClaimSubmissionRequest) (*models.ClaimSubmissionResponse, error) {
	path := "/v1/claims"
	var result models.ClaimSubmissionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RejectClaim Reject claim
func (s *ClaimService) RejectClaim(ctx context.Context, claimId string, req *models.ClaimRejectionRequest) (*models.ClaimRejectionResponse, error) {
	path := "/v1/claims/{claim_id}:reject"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	var result models.ClaimRejectionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RequestMoreDocuments Request more documents from claimant
func (s *ClaimService) RequestMoreDocuments(ctx context.Context, claimId string, req *models.RequestMoreDocumentsRequest) (*models.RequestMoreDocumentsResponse, error) {
	path := "/v1/claims/{claim_id}:request-documents"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	var result models.RequestMoreDocumentsResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ApproveClaim Approve claim
func (s *ClaimService) ApproveClaim(ctx context.Context, claimId string, req *models.ClaimApprovalRequest) (*models.ClaimApprovalResponse, error) {
	path := "/v1/claims/{claim_id}:approve"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	var result models.ClaimApprovalResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserClaims List user claims
func (s *ClaimService) ListUserClaims(ctx context.Context, customerId string) (*models.UserClaimsListingResponse, error) {
	path := "/v1/users/{customer_id}/claims"
	path = strings.ReplaceAll(path, "{customer_id}", customerId)
	var result models.UserClaimsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

