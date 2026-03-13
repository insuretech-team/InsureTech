package models


// PayoutProcessingRequest represents a payout_processing_request
type PayoutProcessingRequest struct {
	PayoutId string `json:"payout_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
}
