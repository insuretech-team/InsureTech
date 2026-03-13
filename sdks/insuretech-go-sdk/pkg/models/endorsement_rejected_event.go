package models

import (
	"time"
)

// EndorsementRejectedEvent represents a endorsement_rejected_event
type EndorsementRejectedEvent struct {
	EventId string `json:"event_id,omitempty"`
	EndorsementId string `json:"endorsement_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
