package models


// WorkflowStartRequest represents a workflow_start_request
type WorkflowStartRequest struct {
	Context map[string]interface{} `json:"context,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	WorkflowDefinitionId string `json:"workflow_definition_id"`
}
