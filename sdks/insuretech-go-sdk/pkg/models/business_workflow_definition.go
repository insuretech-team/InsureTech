package models

import (
	"time"
)

// BusinessWorkflowDefinition represents a business_workflow_definition
type BusinessWorkflowDefinition struct {
	BusinessWorkflowId string `json:"business_workflow_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	InputSchema string `json:"input_schema,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	OutputSchema string `json:"output_schema,omitempty"`
	RulesConfig string `json:"rules_config"`
	Status interface{} `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	Version int `json:"version"`
	WorkflowName string `json:"workflow_name"`
	WorkflowType *BusinessWorkflowType `json:"workflow_type"`
}
