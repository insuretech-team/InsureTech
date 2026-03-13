package models


// TaskTaskCompletionResponse represents a task_task_completion_response
type TaskTaskCompletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
