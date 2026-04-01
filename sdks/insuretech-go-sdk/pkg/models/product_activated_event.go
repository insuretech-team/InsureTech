package models

import (
	"time"
)

// ProductActivatedEvent represents a product_activated_event
type ProductActivatedEvent struct {
	ActivatedBy string `json:"activated_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
