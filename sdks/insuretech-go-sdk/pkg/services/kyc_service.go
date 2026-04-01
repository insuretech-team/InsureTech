package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// KycService handles kyc-related API calls
type KycService struct {
	Client Client
}

// StartKYCVerification Start KYC verification
func (s *KycService) StartKYCVerification(ctx context.Context, req *models.KYCVerificationStartRequest) error {
	path := "/v1/kyc-verifications"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetKYCVerification Get KYC verification
func (s *KycService) GetKYCVerification(ctx context.Context, kycVerificationId string) error {
	path := "/v1/kyc-verifications/{kyc_verification_id}"
	path = strings.ReplaceAll(path, "{kyc_verification_id}", kycVerificationId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UploadDocument Upload document
func (s *KycService) UploadDocument(ctx context.Context, kycVerificationId string, req *models.KycDocumentUploadRequest) error {
	path := "/v1/kyc-verifications/{kyc_verification_id}/documents"
	path = strings.ReplaceAll(path, "{kyc_verification_id}", kycVerificationId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// RejectKYC Reject KYC
func (s *KycService) RejectKYC(ctx context.Context, kycVerificationId string, req *models.KYCRejectionRequest) error {
	path := "/v1/kyc-verifications/{kyc_verification_id}:reject"
	path = strings.ReplaceAll(path, "{kyc_verification_id}", kycVerificationId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// VerifyKYC Verify KYC
func (s *KycService) VerifyKYC(ctx context.Context, kycVerificationId string, req *models.KYCVerificationRequest) error {
	path := "/v1/kyc-verifications/{kyc_verification_id}:verify"
	path = strings.ReplaceAll(path, "{kyc_verification_id}", kycVerificationId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListPendingVerifications List pending KYC verifications (admin review queue)
func (s *KycService) ListPendingVerifications(ctx context.Context) error {
	path := "/v1/kyc-verifications:pending"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

