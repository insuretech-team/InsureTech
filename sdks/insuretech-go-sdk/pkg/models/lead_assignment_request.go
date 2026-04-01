package models


// LeadAssignmentRequest represents a lead_assignment_request
type LeadAssignmentRequest struct {
	AgentId string `json:"agent_id"`
	LeadId string `json:"lead_id"`
}
