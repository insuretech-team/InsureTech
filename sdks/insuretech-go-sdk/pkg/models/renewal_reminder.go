package models

import (
	"time"
)

// RenewalReminder represents a renewal_reminder
type RenewalReminder struct {
	ScheduledAt time.Time `json:"scheduled_at"`
	SentAt time.Time `json:"sent_at,omitempty"`
	NotificationId string `json:"notification_id,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	RenewalScheduleId string `json:"renewal_schedule_id"`
	DaysBeforeRenewal int `json:"days_before_renewal"`
	Channel *ReminderChannel `json:"channel"`
	Status interface{} `json:"status"`
}
