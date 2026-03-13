package models


// EmployeeDeletionResponse represents a employee_deletion_response
type EmployeeDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
