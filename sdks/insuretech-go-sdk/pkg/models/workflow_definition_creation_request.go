package models


// WorkflowDefinitionCreationRequest represents a workflow_definition_creation_request
type WorkflowDefinitionCreationRequest struct {
	Description string `json:"description,omitempty"`
	EntityType string `json:"entity_type"`
	Name string `json:"name"`
	Steps map[string]interface{} `json:"steps,omitempty"`
	Type string `json:"type"`
}
