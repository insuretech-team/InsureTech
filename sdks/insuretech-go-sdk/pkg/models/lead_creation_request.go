package models

import (
	"time"
)

// LeadCreationRequest represents a lead_creation_request
type LeadCreationRequest struct {
	Address string `json:"address,omitempty"`
	AssignedAgentId string `json:"assigned_agent_id"`
	CreatedBy string `json:"created_by,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	DesiredCoverageAmount string `json:"desired_coverage_amount,omitempty"`
	DesiredInsuranceType string `json:"desired_insurance_type,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	LastName string `json:"last_name,omitempty"`
	LeadPriority *LeadPriority `json:"lead_priority,omitempty"`
	LeadSource *LeadSource `json:"lead_source,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	SpecificRequirements string `json:"specific_requirements,omitempty"`
	Title string `json:"title,omitempty"`
}
