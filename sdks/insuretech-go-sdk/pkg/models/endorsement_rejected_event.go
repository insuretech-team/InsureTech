package models

import (
	"time"
)

// EndorsementRejectedEvent represents a endorsement_rejected_event
type EndorsementRejectedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EndorsementId string `json:"endorsement_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
