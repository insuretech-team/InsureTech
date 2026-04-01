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

// SubmitClaim Submit claim
func (s *ClaimService) SubmitClaim(ctx context.Context, req *models.ClaimSubmissionRequest) error {
	path := "/v1/claims"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetClaim Get claim details
func (s *ClaimService) GetClaim(ctx context.Context, claimId string) error {
	path := "/v1/claims/{claim_id}"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UploadDocument Upload claim document
func (s *ClaimService) UploadDocument(ctx context.Context, claimId string, req *models.ClaimsDocumentUploadRequest) error {
	path := "/v1/claims/{claim_id}/documents"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ApproveClaim Approve claim
func (s *ClaimService) ApproveClaim(ctx context.Context, claimId string, req *models.ClaimApprovalRequest) error {
	path := "/v1/claims/{claim_id}:approve"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// DisputeClaim Dispute claim (by customer)
func (s *ClaimService) DisputeClaim(ctx context.Context, claimId string, req *models.DisputeClaimRequest) error {
	path := "/v1/claims/{claim_id}:dispute"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RejectClaim Reject claim
func (s *ClaimService) RejectClaim(ctx context.Context, claimId string, req *models.ClaimRejectionRequest) error {
	path := "/v1/claims/{claim_id}:reject"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RequestMoreDocuments Request more documents from claimant
func (s *ClaimService) RequestMoreDocuments(ctx context.Context, claimId string, req *models.RequestMoreDocumentsRequest) error {
	path := "/v1/claims/{claim_id}:request-documents"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// SettleClaim Settle claim
func (s *ClaimService) SettleClaim(ctx context.Context, claimId string, req *models.SettleClaimRequest) error {
	path := "/v1/claims/{claim_id}:settle"
	path = strings.ReplaceAll(path, "{claim_id}", claimId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListUserClaims List user claims
func (s *ClaimService) ListUserClaims(ctx context.Context, customerId string) error {
	path := "/v1/users/{customer_id}/claims"
	path = strings.ReplaceAll(path, "{customer_id}", customerId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

