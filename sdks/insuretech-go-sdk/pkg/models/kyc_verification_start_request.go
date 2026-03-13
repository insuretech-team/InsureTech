package models


// KYCVerificationStartRequest represents a kyc_verification_start_request
type KYCVerificationStartRequest struct {
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	Method string `json:"method,omitempty"`
	Type string `json:"type"`
}
