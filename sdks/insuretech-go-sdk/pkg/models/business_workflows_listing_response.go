package models


// BusinessWorkflowsListingResponse represents a business_workflows_listing_response
type BusinessWorkflowsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Workflows []*BusinessWorkflowDefinition `json:"workflows,omitempty"`
}
