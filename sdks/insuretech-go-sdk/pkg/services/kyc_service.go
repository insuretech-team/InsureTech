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

// ListPendingVerifications List pending KYC verifications (admin review queue)
func (s *KycService) ListPendingVerifications(ctx context.Context) (*models.PendingVerificationsListingResponse, error) {
	path := "/v1/kyc-verifications:pending"
	var result models.PendingVerificationsListingResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// StartKYCVerification Start KYC verification
func (s *KycService) StartKYCVerification(ctx context.Context, req *models.KYCVerificationStartRequest) (*models.KYCVerificationStartResponse, error) {
	path := "/v1/kyc-verifications"
	var result models.KYCVerificationStartResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadDocument Upload document
func (s *KycService) UploadDocument(ctx context.Context, kycVerificationId string, req *models.KycDocumentUploadRequest) (*models.KycDocumentUploadResponse, error) {
	path := "/v1/kyc-verifications/{kyc_verification_id}/documents"
	path = strings.ReplaceAll(path, "{kyc_verification_id}", kycVerificationId)
	var result models.KycDocumentUploadResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetKYCVerification Get KYC verification
func (s *KycService) GetKYCVerification(ctx context.Context, kycVerificationId string) (*models.KYCVerificationRetrievalResponse, error) {
	path := "/v1/kyc-verifications/{kyc_verification_id}"
	path = strings.ReplaceAll(path, "{kyc_verification_id}", kycVerificationId)
	var result models.KYCVerificationRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VerifyKYC Verify KYC
func (s *KycService) VerifyKYC(ctx context.Context, kycVerificationId string, req *models.KYCVerificationRequest) (*models.KYCVerificationResponse, error) {
	path := "/v1/kyc-verifications/{kyc_verification_id}"
	path = strings.ReplaceAll(path, "{kyc_verification_id}", kycVerificationId)
	var result models.KYCVerificationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

