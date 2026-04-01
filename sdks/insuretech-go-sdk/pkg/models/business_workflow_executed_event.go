package models

import (
	"time"
)

// BusinessWorkflowExecutedEvent represents a business_workflow_executed_event
type BusinessWorkflowExecutedEvent struct {
	BusinessWorkflowId string `json:"business_workflow_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExecutedAt time.Time `json:"executed_at,omitempty"`
	ExecutedBy string `json:"executed_by,omitempty"`
	ExecutionId string `json:"execution_id,omitempty"`
	ExecutionTimeMs int `json:"execution_time_ms,omitempty"`
	IsSuccess bool `json:"is_success,omitempty"`
	RulesFailed int `json:"rules_failed,omitempty"`
	RulesPassed int `json:"rules_passed,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
}
