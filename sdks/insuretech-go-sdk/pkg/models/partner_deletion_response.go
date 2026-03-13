package models


// PartnerDeletionResponse represents a partner_deletion_response
type PartnerDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
