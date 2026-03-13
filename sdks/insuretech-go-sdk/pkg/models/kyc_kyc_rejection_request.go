package models


// KycKYCRejectionRequest represents a kyc_kyc_rejection_request
type KycKYCRejectionRequest struct {
	KycVerificationId string `json:"kyc_verification_id"`
	Reason string `json:"reason,omitempty"`
}
