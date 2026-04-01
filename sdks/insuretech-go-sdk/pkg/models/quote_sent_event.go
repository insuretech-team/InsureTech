package models

import (
	"time"
)

// QuoteSentEvent represents a quote_sent_event
type QuoteSentEvent struct {
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	SentAt time.Time `json:"sent_at,omitempty"`
	SentMethod string `json:"sent_method,omitempty"`
}
