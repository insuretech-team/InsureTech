package models

import (
	"time"
)

// CustomerQueryHandledEvent represents a customer_query_handled_event
type CustomerQueryHandledEvent struct {
	ConversationId string `json:"conversation_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	QueryType string `json:"query_type,omitempty"`
	ResolvedByAi bool `json:"resolved_by_ai,omitempty"`
	ResponseTimeSeconds int `json:"response_time_seconds,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	EscalatedToHuman bool `json:"escalated_to_human,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
