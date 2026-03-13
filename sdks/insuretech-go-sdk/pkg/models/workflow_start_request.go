package models


// WorkflowStartRequest represents a workflow_start_request
type WorkflowStartRequest struct {
	WorkflowDefinitionId string `json:"workflow_definition_id"`
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
	Context map[string]interface{} `json:"context,omitempty"`
}
