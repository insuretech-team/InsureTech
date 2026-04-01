package models

import (
	"time"
)

// LifeProductCreatedEvent represents a life_product_created_event
type LifeProductCreatedEvent struct {
	BaseRate string `json:"base_rate,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	ProductType *LifeProductType `json:"product_type,omitempty"`
}
