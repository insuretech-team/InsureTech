package models


// DepartmentCreationResponse represents a department_creation_response
type DepartmentCreationResponse struct {
	Error *Error `json:"error,omitempty"`
	Department *Department `json:"department,omitempty"`
	Message string `json:"message,omitempty"`
}
