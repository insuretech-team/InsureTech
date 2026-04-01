package models


// CommissionsListingRequest represents a commissions_listing_request
type CommissionsListingRequest struct {
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	RecipientId string `json:"recipient_id"`
	RecipientType string `json:"recipient_type,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	Status string `json:"status,omitempty"`
}
