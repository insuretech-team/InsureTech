package models


// ContactDeletionRequest represents a contact_deletion_request
type ContactDeletionRequest struct {
	ContactId string `json:"contact_id"`
	Permanent bool `json:"permanent,omitempty"`
}
