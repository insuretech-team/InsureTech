package models


// EmployeeCreationResponse represents a employee_creation_response
type EmployeeCreationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	Employee *EmployeeView `json:"employee,omitempty"`
}
