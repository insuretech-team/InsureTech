package models

import (
	"time"
)

// TenantCreatedEvent represents a tenant_created_event
type TenantCreatedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	TenantCode string `json:"tenant_code,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TenantName string `json:"tenant_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
