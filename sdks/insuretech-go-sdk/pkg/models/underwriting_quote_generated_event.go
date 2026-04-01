package models

import (
	"time"
)

// UnderwritingQuoteGeneratedEvent represents a underwriting_quote_generated_event
type UnderwritingQuoteGeneratedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	ValidUntil string `json:"valid_until,omitempty"`
}
