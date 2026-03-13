package models

import (
	"time"
)

// ApiKeyCreatedEvent represents a api_key_created_event
type ApiKeyCreatedEvent struct {
	OwnerId string `json:"owner_id,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ApiKeyId string `json:"api_key_id,omitempty"`
	OwnerType string `json:"owner_type,omitempty"`
}
