package models

import (
	"time"
)

// BeneficiaryStatusChangedEvent represents a beneficiary_status_changed_event
type BeneficiaryStatusChangedEvent struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewStatus string `json:"new_status,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
