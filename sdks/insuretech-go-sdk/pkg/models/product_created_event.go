package models

import (
	"time"
)

// ProductCreatedEvent represents a product_created_event
type ProductCreatedEvent struct {
	ProductName string `json:"product_name,omitempty"`
	Category string `json:"category,omitempty"`
	BasePremium *Money `json:"base_premium,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
}
