package models


// RenewalReminderUpdateResponse represents a renewal_reminder_update_response
type RenewalReminderUpdateResponse struct {
	Reminder *RenewalReminder `json:"reminder,omitempty"`
}
