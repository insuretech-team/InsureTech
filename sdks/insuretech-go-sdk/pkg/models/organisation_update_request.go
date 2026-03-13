package models


// OrganisationUpdateRequest represents a organisation_update_request
type OrganisationUpdateRequest struct {
	OrganisationId string `json:"organisation_id"`
	Name string `json:"name"`
	Industry string `json:"industry,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Address string `json:"address,omitempty"`
	Status *OrganisationStatus `json:"status,omitempty"`
}
