package models

import (
	"time"
)

// ProductDiscontinuedEvent represents a product_discontinued_event
type ProductDiscontinuedEvent struct {
	DiscontinuedBy string `json:"discontinued_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
