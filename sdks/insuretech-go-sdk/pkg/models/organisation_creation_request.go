package models


// OrganisationCreationRequest represents a organisation_creation_request
type OrganisationCreationRequest struct {
	ContactPhone string `json:"contact_phone,omitempty"`
	Address string `json:"address,omitempty"`
	TenantId string `json:"tenant_id"`
	Name string `json:"name"`
	Code string `json:"code,omitempty"`
	Industry string `json:"industry,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
}
