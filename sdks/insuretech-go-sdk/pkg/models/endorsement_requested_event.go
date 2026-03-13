package models

import (
	"time"
)

// EndorsementRequestedEvent represents a endorsement_requested_event
type EndorsementRequestedEvent struct {
	EndorsementNumber string `json:"endorsement_number,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Type string `json:"type,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	EndorsementId string `json:"endorsement_id,omitempty"`
}
