package models


// ContactUpdateRequest represents a contact_update_request
type ContactUpdateRequest struct {
	Address string `json:"address,omitempty"`
	AssignedAgentId string `json:"assigned_agent_id"`
	ContactId string `json:"contact_id"`
	ContactStatus *ContactStatus `json:"contact_status,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	MarketingConsent bool `json:"marketing_consent,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PreferredContactMethod *PreferredContactMethod `json:"preferred_contact_method,omitempty"`
}
