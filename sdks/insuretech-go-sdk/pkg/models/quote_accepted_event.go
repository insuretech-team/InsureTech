package models

import (
	"time"
)

// QuoteAcceptedEvent represents a quote_accepted_event
type QuoteAcceptedEvent struct {
	AcceptedAt time.Time `json:"accepted_at,omitempty"`
	AcceptedPremium *Money `json:"accepted_premium,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
}
