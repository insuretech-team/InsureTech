package models

import (
	"time"
)

// EndorsementApprovedEvent represents a endorsement_approved_event
type EndorsementApprovedEvent struct {
	EventId string `json:"event_id,omitempty"`
	EndorsementId string `json:"endorsement_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Type string `json:"type,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
	EffectiveDate string `json:"effective_date,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
