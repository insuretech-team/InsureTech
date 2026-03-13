package models


// RenewalReminderCreationResponse represents a renewal_reminder_creation_response
type RenewalReminderCreationResponse struct {
	Reminder *RenewalReminder `json:"reminder,omitempty"`
}
