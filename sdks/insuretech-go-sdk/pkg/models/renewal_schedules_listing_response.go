package models


// RenewalSchedulesListingResponse represents a renewal_schedules_listing_response
type RenewalSchedulesListingResponse struct {
	Schedules []*RenewalSchedule `json:"schedules,omitempty"`
}
