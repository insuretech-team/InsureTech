package models

import (
	"time"
)

// Ticket represents a ticket
type Ticket struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	BeneficiaryId string `json:"beneficiary_id"`
	Category *TicketCategory `json:"category"`
	ClosedAt time.Time `json:"closed_at,omitempty"`
	Description string `json:"description"`
	Id string `json:"id"`
	Priority interface{} `json:"priority"`
	RelatedEntityId string `json:"related_entity_id,omitempty"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	ResolvedAt time.Time `json:"resolved_at,omitempty"`
	Status interface{} `json:"status"`
	Subject string `json:"subject"`
	TicketNumber string `json:"ticket_number"`
	Type *TicketType `json:"type"`
}
