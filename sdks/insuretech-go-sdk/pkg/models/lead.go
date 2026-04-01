package models

import (
	"time"
)

// Lead represents a lead
type Lead struct {
	AdditionalEmailAddress string `json:"additional_email_address,omitempty"`
	Address string `json:"address,omitempty"`
	AlternatePhoneNumber string `json:"alternate_phone_number,omitempty"`
	AssignedAgentId string `json:"assigned_agent_id,omitempty"`
	ConversionReason string `json:"conversion_reason,omitempty"`
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	ConvertedContactId string `json:"converted_contact_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DesiredCoverageAmount string `json:"desired_coverage_amount,omitempty"`
	DesiredInsuranceType string `json:"desired_insurance_type,omitempty"`
	EmailAddress string `json:"email_address"`
	FirstName string `json:"first_name"`
	Gender string `json:"gender,omitempty"`
	LastName string `json:"last_name"`
	LeadId string `json:"lead_id"`
	LeadPriority *LeadPriority `json:"lead_priority,omitempty"`
	LeadScore int `json:"lead_score,omitempty"`
	LeadSource *LeadSource `json:"lead_source,omitempty"`
	LeadStatus interface{} `json:"lead_status"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PhoneNumber string `json:"phone_number"`
	QualificationStatus *QualificationStatus `json:"qualification_status,omitempty"`
	SpecificRequirements string `json:"specific_requirements,omitempty"`
	Title string `json:"title,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
