package models


// KycKYCRejectionRequest represents a kyc_kyc_rejection_request
type KycKYCRejectionRequest struct {
	Reason string `json:"reason,omitempty"`
	KycVerificationId string `json:"kyc_verification_id"`
}
