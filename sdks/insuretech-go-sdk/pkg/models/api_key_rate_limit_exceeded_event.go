package models

import (
	"time"
)

// ApiKeyRateLimitExceededEvent represents a api_key_rate_limit_exceeded_event
type ApiKeyRateLimitExceededEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ApiKeyId string `json:"api_key_id,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	RequestsCount int `json:"requests_count,omitempty"`
}
