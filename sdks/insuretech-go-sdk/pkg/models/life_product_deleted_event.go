package models

import (
	"time"
)

// LifeProductDeletedEvent represents a life_product_deleted_event
type LifeProductDeletedEvent struct {
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Permanent bool `json:"permanent,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
}
