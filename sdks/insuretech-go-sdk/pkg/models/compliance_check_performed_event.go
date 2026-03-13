package models

import (
	"time"
)

// ComplianceCheckPerformedEvent represents a compliance_check_performed_event
type ComplianceCheckPerformedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ComplianceLogId string `json:"compliance_log_id,omitempty"`
	Type string `json:"type,omitempty"`
	Regulation string `json:"regulation,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Status string `json:"status,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
}
