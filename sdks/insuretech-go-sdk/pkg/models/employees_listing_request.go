package models


// EmployeesListingRequest represents a employees_listing_request
type EmployeesListingRequest struct {
	BusinessId string `json:"business_id"`
	DepartmentId string `json:"department_id"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	Status *EmployeeStatus `json:"status,omitempty"`
}
