package models

import (
	"time"
)

// LifeProductUpdatedEvent represents a life_product_updated_event
type LifeProductUpdatedEvent struct {
	ChangedFields []string `json:"changed_fields,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
