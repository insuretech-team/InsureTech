package models


// BusinessWorkflowEvaluationRequest represents a business_workflow_evaluation_request
type BusinessWorkflowEvaluationRequest struct {
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	ExecutedBy string `json:"executed_by,omitempty"`
	Inputs map[string]interface{} `json:"inputs,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
}
