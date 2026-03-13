package models


// RenewalRemindersListingResponse represents a renewal_reminders_listing_response
type RenewalRemindersListingResponse struct {
	Reminders []*RenewalReminder `json:"reminders,omitempty"`
}
