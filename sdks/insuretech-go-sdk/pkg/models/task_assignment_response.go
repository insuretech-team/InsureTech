package models


// TaskAssignmentResponse represents a task_assignment_response
type TaskAssignmentResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
