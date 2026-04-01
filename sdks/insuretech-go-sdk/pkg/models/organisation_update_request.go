package models


// OrganisationUpdateRequest represents a organisation_update_request
type OrganisationUpdateRequest struct {
	Address string `json:"address,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Industry string `json:"industry,omitempty"`
	Name string `json:"name"`
	OrganisationId string `json:"organisation_id"`
	Status *OrganisationStatus `json:"status,omitempty"`
}
