package models


// WorkflowDefinitionRetrievalResponse represents a workflow_definition_retrieval_response
type WorkflowDefinitionRetrievalResponse struct {
	WorkflowDefinition *WorkflowDefinition `json:"workflow_definition,omitempty"`
	Error *Error `json:"error,omitempty"`
}
