package models


// TaskUpdateResponse represents a task_update_response
type TaskUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
