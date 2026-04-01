package models


// OrganisationCreationRequest represents a organisation_creation_request
type OrganisationCreationRequest struct {
	Address string `json:"address,omitempty"`
	Code string `json:"code,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Industry string `json:"industry,omitempty"`
	Name string `json:"name"`
	TenantId string `json:"tenant_id"`
}
