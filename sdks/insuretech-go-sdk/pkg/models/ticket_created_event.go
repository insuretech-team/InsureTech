package models

import (
	"time"
)

// TicketCreatedEvent represents a ticket_created_event
type TicketCreatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	TicketId string `json:"ticket_id,omitempty"`
	TicketNumber string `json:"ticket_number,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	Type string `json:"type,omitempty"`
	Category string `json:"category,omitempty"`
	Priority string `json:"priority,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
