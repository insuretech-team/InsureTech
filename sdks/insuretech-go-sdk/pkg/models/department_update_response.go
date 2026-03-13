package models


// DepartmentUpdateResponse represents a department_update_response
type DepartmentUpdateResponse struct {
	Error *Error `json:"error,omitempty"`
	Department *Department `json:"department,omitempty"`
	Message string `json:"message,omitempty"`
}
