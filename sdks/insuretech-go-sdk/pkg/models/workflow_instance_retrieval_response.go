package models


// WorkflowInstanceRetrievalResponse represents a workflow_instance_retrieval_response
type WorkflowInstanceRetrievalResponse struct {
	WorkflowInstance *WorkflowInstance `json:"workflow_instance,omitempty"`
	Tasks []*WorkflowTask `json:"tasks,omitempty"`
	Error *Error `json:"error,omitempty"`
}
