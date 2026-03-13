package models


// DepartmentDeletionResponse represents a department_deletion_response
type DepartmentDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
