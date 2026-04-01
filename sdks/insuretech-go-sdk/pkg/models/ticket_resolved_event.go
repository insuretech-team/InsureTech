package models

import (
	"time"
)

// TicketResolvedEvent represents a ticket_resolved_event
type TicketResolvedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ResolvedBy string `json:"resolved_by,omitempty"`
	TicketId string `json:"ticket_id,omitempty"`
	TicketNumber string `json:"ticket_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
