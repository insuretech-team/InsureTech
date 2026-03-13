package models


// WorkflowTaskCompletionResponse represents a workflow_task_completion_response
type WorkflowTaskCompletionResponse struct {
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
