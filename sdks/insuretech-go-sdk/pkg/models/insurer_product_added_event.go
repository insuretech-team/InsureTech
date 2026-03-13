package models

import (
	"time"
)

// InsurerProductAddedEvent represents a insurer_product_added_event
type InsurerProductAddedEvent struct {
	ProductCode string `json:"product_code,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InsurerProductId string `json:"insurer_product_id,omitempty"`
	InsurerId string `json:"insurer_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
}
