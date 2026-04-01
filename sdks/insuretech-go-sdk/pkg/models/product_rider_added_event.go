package models

import (
	"time"
)

// ProductRiderAddedEvent represents a product_rider_added_event
type ProductRiderAddedEvent struct {
	CoverageAmount *Money `json:"coverage_amount,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IsMandatory bool `json:"is_mandatory,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	RiderId string `json:"rider_id,omitempty"`
	RiderName string `json:"rider_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
