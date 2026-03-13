package models


// RenewalReminderRetrievalResponse represents a renewal_reminder_retrieval_response
type RenewalReminderRetrievalResponse struct {
	Reminder *RenewalReminder `json:"reminder,omitempty"`
}
