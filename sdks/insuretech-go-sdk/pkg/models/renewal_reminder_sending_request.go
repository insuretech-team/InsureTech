package models


// RenewalReminderSendingRequest represents a renewal_reminder_sending_request
type RenewalReminderSendingRequest struct {
	Channel string `json:"channel,omitempty"`
	RenewalScheduleId string `json:"renewal_schedule_id"`
}
