package models

import (
	"time"
)

// ApiKeyCreatedEvent represents a api_key_created_event
type ApiKeyCreatedEvent struct {
	ApiKeyId string `json:"api_key_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	OwnerType string `json:"owner_type,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
