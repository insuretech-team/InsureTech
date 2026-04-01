package models


// BusinessWorkflowCreationRequest represents a business_workflow_creation_request
type BusinessWorkflowCreationRequest struct {
	CreatedBy string `json:"created_by,omitempty"`
	Description string `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	OutputSchema map[string]interface{} `json:"output_schema,omitempty"`
	Rules []*BusinessRule `json:"rules,omitempty"`
	WorkflowName string `json:"workflow_name"`
	WorkflowType *BusinessWorkflowType `json:"workflow_type,omitempty"`
}
