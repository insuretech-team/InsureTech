package models

import (
	"time"
)

// TicketCreatedEvent represents a ticket_created_event
type TicketCreatedEvent struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	Category string `json:"category,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Priority string `json:"priority,omitempty"`
	TicketId string `json:"ticket_id,omitempty"`
	TicketNumber string `json:"ticket_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
