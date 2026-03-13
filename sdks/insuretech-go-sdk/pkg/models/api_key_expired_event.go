package models

import (
	"time"
)

// ApiKeyExpiredEvent represents a api_key_expired_event
type ApiKeyExpiredEvent struct {
	EventId string `json:"event_id,omitempty"`
	ApiKeyId string `json:"api_key_id,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
