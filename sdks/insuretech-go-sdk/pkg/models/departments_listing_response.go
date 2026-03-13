package models


// DepartmentsListingResponse represents a departments_listing_response
type DepartmentsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	Departments []*Department `json:"departments,omitempty"`
}
