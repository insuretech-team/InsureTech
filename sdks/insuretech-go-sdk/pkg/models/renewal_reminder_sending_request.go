package models


// RenewalReminderSendingRequest represents a renewal_reminder_sending_request
type RenewalReminderSendingRequest struct {
	RenewalScheduleId string `json:"renewal_schedule_id"`
	Channel string `json:"channel,omitempty"`
}
