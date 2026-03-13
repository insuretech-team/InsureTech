package models


// EndorsementEndorsementRetrievalResponse represents a endorsement_endorsement_retrieval_response
type EndorsementEndorsementRetrievalResponse struct {
	Endorsement *Endorsement `json:"endorsement,omitempty"`
	Error *Error `json:"error,omitempty"`
}
