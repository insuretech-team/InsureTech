package models

import (
	"time"
)

// InsurerProductStatusChangedEvent represents a insurer_product_status_changed_event
type InsurerProductStatusChangedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InsurerProductId string `json:"insurer_product_id,omitempty"`
	NewStatus string `json:"new_status,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
