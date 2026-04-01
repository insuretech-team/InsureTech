package models

import (
	"time"
)

// TenantStatusChangedEvent represents a tenant_status_changed_event
type TenantStatusChangedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewStatus string `json:"new_status,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
