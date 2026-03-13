package models


// WorkflowTaskCompletionRequest represents a workflow_task_completion_request
type WorkflowTaskCompletionRequest struct {
	TaskId string `json:"task_id"`
	Decision string `json:"decision,omitempty"`
	Comments string `json:"comments,omitempty"`
}
