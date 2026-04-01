package models

import (
	"time"
)

// BusinessWorkflowExecution represents a business_workflow_execution
type BusinessWorkflowExecution struct {
	BusinessWorkflowId string `json:"business_workflow_id"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	ExecutedAt time.Time `json:"executed_at"`
	ExecutedBy string `json:"executed_by,omitempty"`
	ExecutionId string `json:"execution_id"`
	ExecutionTimeMs int `json:"execution_time_ms,omitempty"`
	InputsJson string `json:"inputs_json,omitempty"`
	IsSuccess bool `json:"is_success"`
	OutputsJson string `json:"outputs_json,omitempty"`
	Results []*BusinessRuleResult `json:"results,omitempty"`
}
