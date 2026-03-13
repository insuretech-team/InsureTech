package models


// UpcomingRenewalsListingResponse represents a upcoming_renewals_listing_response
type UpcomingRenewalsListingResponse struct {
	RenewalSchedules []*RenewalSchedule `json:"renewal_schedules,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
