package models


// OrganisationsListingRequest represents a organisations_listing_request
type OrganisationsListingRequest struct {
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	Status *OrganisationStatus `json:"status,omitempty"`
	TenantId string `json:"tenant_id"`
}
