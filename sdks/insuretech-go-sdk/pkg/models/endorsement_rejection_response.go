package models


// EndorsementRejectionResponse represents a endorsement_rejection_response
type EndorsementRejectionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
