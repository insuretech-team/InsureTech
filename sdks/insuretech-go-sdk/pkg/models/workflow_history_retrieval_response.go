package models


// WorkflowHistoryRetrievalResponse represents a workflow_history_retrieval_response
type WorkflowHistoryRetrievalResponse struct {
	WorkflowInstances []*WorkflowInstance `json:"workflow_instances,omitempty"`
	Error *Error `json:"error,omitempty"`
}
