package models


// CommissionsListingRequest represents a commissions_listing_request
type CommissionsListingRequest struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	RecipientType string `json:"recipient_type,omitempty"`
	RecipientId string `json:"recipient_id"`
	Status string `json:"status,omitempty"`
}
