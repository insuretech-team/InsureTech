package models


// RefundProcessingRequest represents a refund_processing_request
type RefundProcessingRequest struct {
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	RefundId string `json:"refund_id"`
}
