package models

import (
	"time"
)

// ApiKeyRevokedEvent represents a api_key_revoked_event
type ApiKeyRevokedEvent struct {
	ApiKeyId string `json:"api_key_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
