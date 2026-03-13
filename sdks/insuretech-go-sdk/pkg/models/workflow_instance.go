package models

import (
	"time"
)

// WorkflowInstance represents a workflow_instance
type WorkflowInstance struct {
	CompletedAt time.Time `json:"completed_at,omitempty"`
	WorkflowDefinitionId string `json:"workflow_definition_id"`
	EntityId string `json:"entity_id"`
	CurrentStep string `json:"current_step,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	EntityType string `json:"entity_type"`
	Status interface{} `json:"status"`
	Context string `json:"context,omitempty"`
	InitiatedBy string `json:"initiated_by"`
	StartedAt time.Time `json:"started_at"`
}
