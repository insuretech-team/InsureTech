package models

import (
	"time"
)

// ConversationStartedEvent represents a conversation_started_event
type ConversationStartedEvent struct {
	UserId string `json:"user_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
	Channel string `json:"channel,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ConversationId string `json:"conversation_id,omitempty"`
}
