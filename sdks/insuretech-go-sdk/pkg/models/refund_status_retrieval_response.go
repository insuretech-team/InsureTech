package models

import (
	"time"
)

// RefundStatusRetrievalResponse represents a refund_status_retrieval_response
type RefundStatusRetrievalResponse struct {
	CompletedAt time.Time `json:"completed_at,omitempty"`
	RefundAmount *Money `json:"refund_amount,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	Status string `json:"status,omitempty"`
}
