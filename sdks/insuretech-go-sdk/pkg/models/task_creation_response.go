package models


// TaskCreationResponse represents a task_creation_response
type TaskCreationResponse struct {
	Error *Error `json:"error,omitempty"`
	TaskId string `json:"task_id,omitempty"`
	Message string `json:"message,omitempty"`
}
