package models


// WorkflowInstanceRetrievalResponse represents a workflow_instance_retrieval_response
type WorkflowInstanceRetrievalResponse struct {
	Tasks []*WorkflowTask `json:"tasks,omitempty"`
	WorkflowInstance *WorkflowInstance `json:"workflow_instance,omitempty"`
}
