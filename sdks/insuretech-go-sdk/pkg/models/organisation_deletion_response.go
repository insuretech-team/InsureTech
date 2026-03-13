package models


// OrganisationDeletionResponse represents a organisation_deletion_response
type OrganisationDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
