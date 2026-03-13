package models

import (
	"time"
)

// AgentRegisteredEvent represents a agent_registered_event
type AgentRegisteredEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
