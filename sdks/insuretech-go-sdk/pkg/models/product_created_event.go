package models

import (
	"time"
)

// ProductCreatedEvent represents a product_created_event
type ProductCreatedEvent struct {
	BasePremium *Money `json:"base_premium,omitempty"`
	Category string `json:"category,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
