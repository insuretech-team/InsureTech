package models


// WorkflowDefinition represents a workflow_definition
type WorkflowDefinition struct {
	AuditInfo interface{} `json:"audit_info"`
	Conditions string `json:"conditions,omitempty"`
	Description string `json:"description,omitempty"`
	EntityType string `json:"entity_type"`
	Name string `json:"name"`
	Status interface{} `json:"status"`
	Steps string `json:"steps"`
	Type *WorkflowType `json:"type"`
	Version int `json:"version"`
	WorkflowDefinitionId string `json:"workflow_definition_id"`
}
