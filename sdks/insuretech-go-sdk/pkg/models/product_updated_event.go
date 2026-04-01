package models

import (
	"time"
)

// ProductUpdatedEvent represents a product_updated_event
type ProductUpdatedEvent struct {
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
