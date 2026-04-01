package models

import (
	"time"
)

// ProductDeactivatedEvent represents a product_deactivated_event
type ProductDeactivatedEvent struct {
	DeactivatedBy string `json:"deactivated_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
