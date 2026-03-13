package models

import (
	"time"
)

// TicketResolvedEvent represents a ticket_resolved_event
type TicketResolvedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	TicketId string `json:"ticket_id,omitempty"`
	TicketNumber string `json:"ticket_number,omitempty"`
	ResolvedBy string `json:"resolved_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
