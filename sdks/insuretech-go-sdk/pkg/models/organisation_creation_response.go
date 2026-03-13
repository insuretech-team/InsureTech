package models


// OrganisationCreationResponse represents a organisation_creation_response
type OrganisationCreationResponse struct {
	Organisation *Organisation `json:"organisation,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
