package models

import (
	"time"
)

// InitiateRefundResponse represents a initiate_refund_response
type InitiateRefundResponse struct {
	InitiatedAt time.Time `json:"initiated_at,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	Status string `json:"status,omitempty"`
}
