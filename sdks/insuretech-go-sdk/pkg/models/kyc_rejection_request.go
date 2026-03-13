package models


// KYCRejectionRequest represents a kyc_rejection_request
type KYCRejectionRequest struct {
	KycVerificationId string `json:"kyc_verification_id"`
	Reason string `json:"reason,omitempty"`
}
