package models

import (
	"time"
)

// ContactCreationRequest represents a contact_creation_request
type ContactCreationRequest struct {
	Address string `json:"address,omitempty"`
	AssignedAgentId string `json:"assigned_agent_id"`
	ContactType *ContactType `json:"contact_type,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	LastName string `json:"last_name,omitempty"`
	MarketingConsent bool `json:"marketing_consent,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PreferredContactMethod *PreferredContactMethod `json:"preferred_contact_method,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	ReferralSource string `json:"referral_source,omitempty"`
}
