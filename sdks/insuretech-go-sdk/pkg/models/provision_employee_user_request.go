package models


// ProvisionEmployeeUserRequest represents a provision_employee_user_request
type ProvisionEmployeeUserRequest struct {
	BusinessId string `json:"business_id"`
	Email string `json:"email"`
	EmployeeId string `json:"employee_id"`
	FullName string `json:"full_name,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
}
