package models


// BusinessWorkflowsListingRequest represents a business_workflows_listing_request
type BusinessWorkflowsListingRequest struct {
	Filter string `json:"filter,omitempty"`
	OrderBy string `json:"order_by,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	Status *BusinessWorkflowStatus `json:"status,omitempty"`
	WorkflowType *BusinessWorkflowType `json:"workflow_type"`
}
