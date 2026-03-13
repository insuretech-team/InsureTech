package models

import (
	"time"
)

// AIDecisionMadeEvent represents a ai_decision_made_event
type AIDecisionMadeEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	Reasoning []string `json:"reasoning,omitempty"`
	DecisionType string `json:"decision_type,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	Decision string `json:"decision,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
