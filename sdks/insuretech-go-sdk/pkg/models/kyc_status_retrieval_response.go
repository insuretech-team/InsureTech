package models

import (
	"time"
)

// KYCStatusRetrievalResponse represents a kyc_status_retrieval_response
type KYCStatusRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	KycId string `json:"kyc_id,omitempty"`
	Status string `json:"status,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	SubmittedAt time.Time `json:"submitted_at,omitempty"`
	ReviewedAt time.Time `json:"reviewed_at,omitempty"`
}
