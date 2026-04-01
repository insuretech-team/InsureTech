package models


// TaskAssignmentRequest represents a task_assignment_request
type TaskAssignmentRequest struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	TaskId string `json:"task_id"`
}
