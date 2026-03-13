package models


// OrganisationsListingRequest represents a organisations_listing_request
type OrganisationsListingRequest struct {
	PageToken string `json:"page_token,omitempty"`
	TenantId string `json:"tenant_id"`
	Status *OrganisationStatus `json:"status,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
