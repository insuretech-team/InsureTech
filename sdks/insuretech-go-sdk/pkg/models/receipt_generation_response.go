package models


// ReceiptGenerationResponse represents a receipt_generation_response
type ReceiptGenerationResponse struct {
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	Error *Error `json:"error,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
}
