package models

import (
	"time"
)

// AIAgentActivatedEvent represents a ai_agent_activated_event
type AIAgentActivatedEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
