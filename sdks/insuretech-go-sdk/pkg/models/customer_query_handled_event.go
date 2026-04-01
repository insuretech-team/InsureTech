package models

import (
	"time"
)

// CustomerQueryHandledEvent represents a customer_query_handled_event
type CustomerQueryHandledEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	ConversationId string `json:"conversation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EscalatedToHuman bool `json:"escalated_to_human,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QueryType string `json:"query_type,omitempty"`
	ResolvedByAi bool `json:"resolved_by_ai,omitempty"`
	ResponseTimeSeconds int `json:"response_time_seconds,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
