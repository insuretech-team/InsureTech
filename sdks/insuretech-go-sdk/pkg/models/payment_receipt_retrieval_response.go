package models

import (
	"time"
)

// PaymentReceiptRetrievalResponse represents a payment_receipt_retrieval_response
type PaymentReceiptRetrievalResponse struct {
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	ReceiptUrl string `json:"receipt_url,omitempty"`
}
