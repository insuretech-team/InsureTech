package models


// EmployeeLoginOrganisationsListingResponse represents a employee_login_organisations_listing_response
type EmployeeLoginOrganisationsListingResponse struct {
	Organisations []*EmployeeLoginOrganisation `json:"organisations,omitempty"`
}
