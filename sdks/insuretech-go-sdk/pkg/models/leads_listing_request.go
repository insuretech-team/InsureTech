package models


// LeadsListingRequest represents a leads_listing_request
type LeadsListingRequest struct {
	AssignedAgentId string `json:"assigned_agent_id"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	SearchQuery string `json:"search_query,omitempty"`
	Source *LeadSource `json:"source,omitempty"`
	Status *LeadStatus `json:"status,omitempty"`
}
