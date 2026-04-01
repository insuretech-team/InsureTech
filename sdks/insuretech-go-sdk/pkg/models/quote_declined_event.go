package models

import (
	"time"
)

// QuoteDeclinedEvent represents a quote_declined_event
type QuoteDeclinedEvent struct {
	CustomerId string `json:"customer_id,omitempty"`
	DeclineReason string `json:"decline_reason,omitempty"`
	DeclinedAt time.Time `json:"declined_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
}
