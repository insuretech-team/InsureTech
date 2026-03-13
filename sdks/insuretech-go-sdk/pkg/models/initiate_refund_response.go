package models

import (
	"time"
)

// InitiateRefundResponse represents a initiate_refund_response
type InitiateRefundResponse struct {
	RefundId string `json:"refund_id,omitempty"`
	Status string `json:"status,omitempty"`
	InitiatedAt time.Time `json:"initiated_at,omitempty"`
	Error *Error `json:"error,omitempty"`
}
