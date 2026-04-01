package models

import (
	"time"
)

// QuoteDeletedEvent represents a quote_deleted_event
type QuoteDeletedEvent struct {
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Permanent bool `json:"permanent,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
}
