package models

import (
	"time"
)

// ProductRiderRemovedEvent represents a product_rider_removed_event
type ProductRiderRemovedEvent struct {
	EventId string `json:"event_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	RemovedBy string `json:"removed_by,omitempty"`
	RiderId string `json:"rider_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
