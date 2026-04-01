package models


// BusinessWorkflowUpdateRequest represents a business_workflow_update_request
type BusinessWorkflowUpdateRequest struct {
	BusinessWorkflowId string `json:"business_workflow_id"`
	Description string `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	OutputSchema map[string]interface{} `json:"output_schema,omitempty"`
	Rules []*BusinessRule `json:"rules,omitempty"`
	Status *BusinessWorkflowStatus `json:"status,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
	WorkflowType *BusinessWorkflowType `json:"workflow_type,omitempty"`
}
