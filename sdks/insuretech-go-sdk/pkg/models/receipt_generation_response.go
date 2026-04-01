package models


// ReceiptGenerationResponse represents a receipt_generation_response
type ReceiptGenerationResponse struct {
	PaymentId string `json:"payment_id,omitempty"`
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
}
