package models


// TaskCompletionRequest represents a task_completion_request
type TaskCompletionRequest struct {
	Comments string `json:"comments,omitempty"`
	TaskId string `json:"task_id"`
}
