package models


// EmployeesListingRequest represents a employees_listing_request
type EmployeesListingRequest struct {
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	DepartmentId string `json:"department_id"`
	BusinessId string `json:"business_id"`
	Status *EmployeeStatus `json:"status,omitempty"`
}
