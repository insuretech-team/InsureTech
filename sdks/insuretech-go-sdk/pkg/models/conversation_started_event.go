package models

import (
	"time"
)

// ConversationStartedEvent represents a conversation_started_event
type ConversationStartedEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
	Channel string `json:"channel,omitempty"`
	ConversationId string `json:"conversation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
