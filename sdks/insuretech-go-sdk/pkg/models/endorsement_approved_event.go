package models

import (
	"time"
)

// EndorsementApprovedEvent represents a endorsement_approved_event
type EndorsementApprovedEvent struct {
	ApprovedBy string `json:"approved_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EffectiveDate string `json:"effective_date,omitempty"`
	EndorsementId string `json:"endorsement_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
