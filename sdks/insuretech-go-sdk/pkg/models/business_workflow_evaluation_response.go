package models

import (
	"time"
)

// BusinessWorkflowEvaluationResponse represents a business_workflow_evaluation_response
type BusinessWorkflowEvaluationResponse struct {
	ExecutedAt time.Time `json:"executed_at,omitempty"`
	ExecutionId string `json:"execution_id,omitempty"`
	ExecutionTimeMs int `json:"execution_time_ms,omitempty"`
	IsSuccess bool `json:"is_success,omitempty"`
	Outputs map[string]interface{} `json:"outputs,omitempty"`
	Results []*BusinessRuleResult `json:"results,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
}
