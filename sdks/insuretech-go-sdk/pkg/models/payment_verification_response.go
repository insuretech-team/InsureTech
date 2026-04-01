package models


// PaymentVerificationResponse represents a payment_verification_response
type PaymentVerificationResponse struct {
	Payment *Payment `json:"payment,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	Status string `json:"status,omitempty"`
	Verified bool `json:"verified,omitempty"`
}
