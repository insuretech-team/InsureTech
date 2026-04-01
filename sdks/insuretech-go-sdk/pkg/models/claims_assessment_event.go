package models

import (
	"time"
)

// ClaimsAssessmentEvent represents a claims_assessment_event
type ClaimsAssessmentEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	AssessmentResult string `json:"assessment_result,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Reasons []string `json:"reasons,omitempty"`
	RecommendedAmount *Money `json:"recommended_amount,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
