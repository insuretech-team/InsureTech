package models


// OrganisationUpdateResponse represents a organisation_update_response
type OrganisationUpdateResponse struct {
	Organisation *Organisation `json:"organisation,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
