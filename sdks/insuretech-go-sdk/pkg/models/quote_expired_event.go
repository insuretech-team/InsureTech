package models

import (
	"time"
)

// QuoteExpiredEvent represents a quote_expired_event
type QuoteExpiredEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
