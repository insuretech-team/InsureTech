package models


// EmployeeUpdateResponse represents a employee_update_response
type EmployeeUpdateResponse struct {
	Employee *EmployeeView `json:"employee,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
