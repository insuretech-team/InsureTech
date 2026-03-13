package models

import (
	"time"
)

// PaymentReceiptRetrievalResponse represents a payment_receipt_retrieval_response
type PaymentReceiptRetrievalResponse struct {
	PaymentId string `json:"payment_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	ReceiptUrl string `json:"receipt_url,omitempty"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	Error *Error `json:"error,omitempty"`
}
