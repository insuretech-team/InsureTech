package models


// LeadUpdateRequest represents a lead_update_request
type LeadUpdateRequest struct {
	Address string `json:"address,omitempty"`
	AssignedAgentId string `json:"assigned_agent_id"`
	DesiredCoverageAmount string `json:"desired_coverage_amount,omitempty"`
	DesiredInsuranceType string `json:"desired_insurance_type,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	LeadId string `json:"lead_id"`
	LeadPriority *LeadPriority `json:"lead_priority,omitempty"`
	LeadScore int `json:"lead_score,omitempty"`
	LeadStatus *LeadStatus `json:"lead_status,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	QualificationStatus *QualificationStatus `json:"qualification_status,omitempty"`
	Title string `json:"title,omitempty"`
}
