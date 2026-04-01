package models


// PartnerVerificationRequest represents a partner_verification_request
type PartnerVerificationRequest struct {
	PartnerId string `json:"partner_id"`
	VerificationData map[string]interface{} `json:"verification_data,omitempty"`
	VerificationType string `json:"verification_type,omitempty"`
}
