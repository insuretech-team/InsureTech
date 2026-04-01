package models

import (
	"time"
)

// Contact represents a contact
type Contact struct {
	AdditionalEmailAddress string `json:"additional_email_address,omitempty"`
	Address string `json:"address,omitempty"`
	AlternatePhoneNumber string `json:"alternate_phone_number,omitempty"`
	AssignedAgentId string `json:"assigned_agent_id,omitempty"`
	ContactId string `json:"contact_id"`
	ContactStatus interface{} `json:"contact_status"`
	ContactType interface{} `json:"contact_type"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	EmailAddress string `json:"email_address"`
	FirstName string `json:"first_name"`
	Gender string `json:"gender,omitempty"`
	LastName string `json:"last_name"`
	MarketingConsent bool `json:"marketing_consent,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PhoneNumber string `json:"phone_number"`
	PreferredContactMethod interface{} `json:"preferred_contact_method,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	ReferralSource string `json:"referral_source,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
