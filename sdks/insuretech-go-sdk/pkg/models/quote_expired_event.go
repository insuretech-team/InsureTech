package models

import (
	"time"
)

// QuoteExpiredEvent represents a quote_expired_event
type QuoteExpiredEvent struct {
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
