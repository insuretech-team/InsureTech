package models

import (
	"time"
)

// AIAgentErrorEvent represents a ai_agent_error_event
type AIAgentErrorEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	ErrorType string `json:"error_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Context string `json:"context,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
