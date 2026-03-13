package models

import (
	"time"
)

// RefundStatusRetrievalResponse represents a refund_status_retrieval_response
type RefundStatusRetrievalResponse struct {
	RefundId string `json:"refund_id,omitempty"`
	Status string `json:"status,omitempty"`
	RefundAmount *Money `json:"refund_amount,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Error *Error `json:"error,omitempty"`
}
