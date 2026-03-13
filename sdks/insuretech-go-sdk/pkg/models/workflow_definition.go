package models


// WorkflowDefinition represents a workflow_definition
type WorkflowDefinition struct {
	WorkflowDefinitionId string `json:"workflow_definition_id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	EntityType string `json:"entity_type"`
	Steps string `json:"steps"`
	Conditions string `json:"conditions,omitempty"`
	Version int `json:"version"`
	Status interface{} `json:"status"`
	Type *WorkflowType `json:"type"`
	AuditInfo interface{} `json:"audit_info"`
}
