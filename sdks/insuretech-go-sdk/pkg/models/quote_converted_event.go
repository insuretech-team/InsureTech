package models

import (
	"time"
)

// QuoteConvertedEvent represents a quote_converted_event
type QuoteConvertedEvent struct {
	PolicyId string `json:"policy_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
}
