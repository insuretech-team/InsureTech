package models


// LeadDeletionRequest represents a lead_deletion_request
type LeadDeletionRequest struct {
	LeadId string `json:"lead_id"`
	Permanent bool `json:"permanent,omitempty"`
}
