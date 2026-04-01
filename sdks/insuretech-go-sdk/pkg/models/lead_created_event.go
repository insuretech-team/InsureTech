package models

import (
	"time"
)

// LeadCreatedEvent represents a lead_created_event
type LeadCreatedEvent struct {
	AssignedAgentId string `json:"assigned_agent_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	LeadId string `json:"lead_id,omitempty"`
	LeadSource *LeadSource `json:"lead_source,omitempty"`
}
