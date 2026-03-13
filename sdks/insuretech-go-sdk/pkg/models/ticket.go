package models

import (
	"time"
)

// Ticket represents a ticket
type Ticket struct {
	ClosedAt time.Time `json:"closed_at,omitempty"`
	TicketNumber string `json:"ticket_number"`
	Type *TicketType `json:"type"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	ResolvedAt time.Time `json:"resolved_at,omitempty"`
	Category *TicketCategory `json:"category"`
	Priority interface{} `json:"priority"`
	Description string `json:"description"`
	RelatedEntityId string `json:"related_entity_id,omitempty"`
	AssignedTo string `json:"assigned_to,omitempty"`
	Id string `json:"id"`
	BeneficiaryId string `json:"beneficiary_id"`
	Subject string `json:"subject"`
	Status interface{} `json:"status"`
	AuditInfo interface{} `json:"audit_info"`
}
