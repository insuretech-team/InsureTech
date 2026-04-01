package models


// UpcomingRenewalsListingRequest represents a upcoming_renewals_listing_request
type UpcomingRenewalsListingRequest struct {
	DaysAhead int `json:"days_ahead"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
}
