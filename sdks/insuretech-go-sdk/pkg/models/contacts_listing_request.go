package models


// ContactsListingRequest represents a contacts_listing_request
type ContactsListingRequest struct {
	AssignedAgentId string `json:"assigned_agent_id"`
	ContactType *ContactType `json:"contact_type,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	SearchQuery string `json:"search_query,omitempty"`
	Status *ContactStatus `json:"status,omitempty"`
}
