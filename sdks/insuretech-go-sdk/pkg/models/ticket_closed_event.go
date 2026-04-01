package models

import (
	"time"
)

// TicketClosedEvent represents a ticket_closed_event
type TicketClosedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	TicketId string `json:"ticket_id,omitempty"`
	TicketNumber string `json:"ticket_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
