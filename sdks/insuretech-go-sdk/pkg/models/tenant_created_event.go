package models

import (
	"time"
)

// TenantCreatedEvent represents a tenant_created_event
type TenantCreatedEvent struct {
	Type string `json:"type,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TenantCode string `json:"tenant_code,omitempty"`
	TenantName string `json:"tenant_name,omitempty"`
}
