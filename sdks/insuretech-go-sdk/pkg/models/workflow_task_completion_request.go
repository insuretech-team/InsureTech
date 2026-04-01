package models


// WorkflowTaskCompletionRequest represents a workflow_task_completion_request
type WorkflowTaskCompletionRequest struct {
	Comments string `json:"comments,omitempty"`
	Decision string `json:"decision,omitempty"`
	TaskId string `json:"task_id"`
}
