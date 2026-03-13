package models


// ReceiptGenerationRequest represents a receipt_generation_request
type ReceiptGenerationRequest struct {
	PaymentId string `json:"payment_id"`
	RequestedBy string `json:"requested_by,omitempty"`
}
