package models

import (
	"time"
)

// InsurerProductStatusChangedEvent represents a insurer_product_status_changed_event
type InsurerProductStatusChangedEvent struct {
	NewStatus string `json:"new_status,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InsurerProductId string `json:"insurer_product_id,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
}
