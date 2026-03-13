package models


// KYCVerificationRequest represents a kyc_verification_request
type KYCVerificationRequest struct {
	KycVerificationId string `json:"kyc_verification_id"`
	VerifiedBy string `json:"verified_by,omitempty"`
	VerificationResult string `json:"verification_result,omitempty"`
}
