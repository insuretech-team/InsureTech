package models

import (
	"time"
)

// ComplianceCheckPerformedEvent represents a compliance_check_performed_event
type ComplianceCheckPerformedEvent struct {
	ComplianceLogId string `json:"compliance_log_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Regulation string `json:"regulation,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
