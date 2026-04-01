package models

import (
	"time"
)

// LeadAssignedEvent represents a lead_assigned_event
type LeadAssignedEvent struct {
	AssignedAt time.Time `json:"assigned_at,omitempty"`
	AssignedBy string `json:"assigned_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	LeadId string `json:"lead_id,omitempty"`
	LeadName string `json:"lead_name,omitempty"`
	NewAgentId string `json:"new_agent_id,omitempty"`
	PreviousAgentId string `json:"previous_agent_id,omitempty"`
}
