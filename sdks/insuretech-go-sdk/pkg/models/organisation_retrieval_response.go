package models


// OrganisationRetrievalResponse represents a organisation_retrieval_response
type OrganisationRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	Organisation *Organisation `json:"organisation,omitempty"`
}
