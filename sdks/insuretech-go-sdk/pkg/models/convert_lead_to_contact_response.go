package models


// ConvertLeadToContactResponse represents a convert_lead_to_contact_response
type ConvertLeadToContactResponse struct {
	Contact *Contact `json:"contact,omitempty"`
	Lead *Lead `json:"lead,omitempty"`
}
