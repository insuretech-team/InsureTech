package models


// EmployeeLoginOrganisationsListingRequest represents a employee_login_organisations_listing_request
type EmployeeLoginOrganisationsListingRequest struct {
	PageSize int `json:"page_size,omitempty"`
	Query string `json:"query"`
}
