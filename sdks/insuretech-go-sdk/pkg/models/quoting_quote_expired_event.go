package models

import (
	"time"
)

// QuotingQuoteExpiredEvent represents a quoting_quote_expired_event
type QuotingQuoteExpiredEvent struct {
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
}
