package models

import (
	"time"
)

// BusinessRuleFailedEvent represents a business_rule_failed_event
type BusinessRuleFailedEvent struct {
	BusinessWorkflowId string `json:"business_workflow_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExecutionId string `json:"execution_id,omitempty"`
	FailedAt time.Time `json:"failed_at,omitempty"`
	RuleName string `json:"rule_name,omitempty"`
}
