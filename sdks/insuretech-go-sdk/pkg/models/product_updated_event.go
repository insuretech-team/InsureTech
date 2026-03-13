package models

import (
	"time"
)

// ProductUpdatedEvent represents a product_updated_event
type ProductUpdatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
