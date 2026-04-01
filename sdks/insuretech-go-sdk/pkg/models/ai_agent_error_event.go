package models

import (
	"time"
)

// AIAgentErrorEvent represents a ai_agent_error_event
type AIAgentErrorEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	Context string `json:"context,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	ErrorType string `json:"error_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
