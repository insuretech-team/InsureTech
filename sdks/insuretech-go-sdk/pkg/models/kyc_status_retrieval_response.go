package models

import (
	"time"
)

// KYCStatusRetrievalResponse represents a kyc_status_retrieval_response
type KYCStatusRetrievalResponse struct {
	CompletedSteps int `json:"completed_steps,omitempty"`
	CurrentStep *KYCStep `json:"current_step,omitempty"`
	KycId string `json:"kyc_id,omitempty"`
	OverallProgress float64 `json:"overall_progress,omitempty"`
	Provider string `json:"provider,omitempty"`
	ProviderReference string `json:"provider_reference,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	RemainingSeconds int `json:"remaining_seconds,omitempty"`
	ReviewedAt time.Time `json:"reviewed_at,omitempty"`
	SessionState string `json:"session_state,omitempty"`
	Status string `json:"status,omitempty"`
	SubmittedAt time.Time `json:"submitted_at,omitempty"`
	TotalSteps int `json:"total_steps,omitempty"`
}
