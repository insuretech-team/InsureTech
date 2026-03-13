package models


// EndorsementRetrievalResponse represents a endorsement_retrieval_response
type EndorsementRetrievalResponse struct {
	Endorsement *Endorsement `json:"endorsement,omitempty"`
	Error *Error `json:"error,omitempty"`
}
