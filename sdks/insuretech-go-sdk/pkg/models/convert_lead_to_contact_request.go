package models


// ConvertLeadToContactRequest represents a convert_lead_to_contact_request
type ConvertLeadToContactRequest struct {
	ConversionReason string `json:"conversion_reason,omitempty"`
	LeadId string `json:"lead_id"`
}
