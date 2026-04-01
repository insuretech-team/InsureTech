package models


// BusinessWorkflowRetrievalResponse represents a business_workflow_retrieval_response
type BusinessWorkflowRetrievalResponse struct {
	Workflow *BusinessWorkflowDefinition `json:"workflow,omitempty"`
}
