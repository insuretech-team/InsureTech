package models

import (
	"time"
)

// QuoteViewedEvent represents a quote_viewed_event
type QuoteViewedEvent struct {
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	ViewedAt time.Time `json:"viewed_at,omitempty"`
}
