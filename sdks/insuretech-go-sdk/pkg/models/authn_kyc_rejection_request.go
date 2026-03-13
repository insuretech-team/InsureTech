package models


// AuthnKYCRejectionRequest represents a authn_kyc_rejection_request
type AuthnKYCRejectionRequest struct {
	RejectionReason string `json:"rejection_reason,omitempty"`
	KycId string `json:"kyc_id"`
	ReviewerId string `json:"reviewer_id"`
}
