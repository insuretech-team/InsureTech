package models


// EndorsementRejectionRequest represents a endorsement_rejection_request
type EndorsementRejectionRequest struct {
	EndorsementId string `json:"endorsement_id"`
	Reason string `json:"reason,omitempty"`
}
