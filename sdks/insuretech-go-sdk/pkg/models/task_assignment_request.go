package models


// TaskAssignmentRequest represents a task_assignment_request
type TaskAssignmentRequest struct {
	TaskId string `json:"task_id"`
	AssignedTo string `json:"assigned_to,omitempty"`
}
