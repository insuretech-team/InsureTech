package models

import (
	"time"
)

// QuoteConvertedEvent represents a quote_converted_event
type QuoteConvertedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
