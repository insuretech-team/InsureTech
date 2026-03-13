package models


// EmployeesListingResponse represents a employees_listing_response
type EmployeesListingResponse struct {
	Employees []*EmployeeView `json:"employees,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
