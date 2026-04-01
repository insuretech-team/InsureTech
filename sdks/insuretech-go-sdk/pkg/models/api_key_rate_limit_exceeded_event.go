package models

import (
	"time"
)

// ApiKeyRateLimitExceededEvent represents a api_key_rate_limit_exceeded_event
type ApiKeyRateLimitExceededEvent struct {
	ApiKeyId string `json:"api_key_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	RequestsCount int `json:"requests_count,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
