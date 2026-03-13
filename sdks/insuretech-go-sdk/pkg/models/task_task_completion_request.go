package models


// TaskTaskCompletionRequest represents a task_task_completion_request
type TaskTaskCompletionRequest struct {
	TaskId string `json:"task_id"`
	Comments string `json:"comments,omitempty"`
}
