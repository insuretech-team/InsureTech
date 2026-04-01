package models

import (
	"time"
)

// EndorsementRequestedEvent represents a endorsement_requested_event
type EndorsementRequestedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EndorsementId string `json:"endorsement_id,omitempty"`
	EndorsementNumber string `json:"endorsement_number,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
