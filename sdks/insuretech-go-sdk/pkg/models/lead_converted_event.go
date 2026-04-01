package models

import (
	"time"
)

// LeadConvertedEvent represents a lead_converted_event
type LeadConvertedEvent struct {
	AssignedAgentId string `json:"assigned_agent_id,omitempty"`
	ContactId string `json:"contact_id,omitempty"`
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	LeadId string `json:"lead_id,omitempty"`
}
