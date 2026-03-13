package models


// UpcomingRenewalsListingRequest represents a upcoming_renewals_listing_request
type UpcomingRenewalsListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	DaysAhead int `json:"days_ahead"`
	Status string `json:"status,omitempty"`
}
