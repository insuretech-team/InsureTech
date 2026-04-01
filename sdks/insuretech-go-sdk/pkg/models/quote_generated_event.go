package models

import (
	"time"
)

// QuoteGeneratedEvent represents a quote_generated_event
type QuoteGeneratedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	ValidUntil string `json:"valid_until,omitempty"`
}
