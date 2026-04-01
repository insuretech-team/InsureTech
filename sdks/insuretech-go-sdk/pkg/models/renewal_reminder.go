package models

import (
	"time"
)

// RenewalReminder represents a renewal_reminder
type RenewalReminder struct {
	AuditInfo interface{} `json:"audit_info"`
	Channel *ReminderChannel `json:"channel"`
	DaysBeforeRenewal int `json:"days_before_renewal"`
	Id string `json:"id"`
	NotificationId string `json:"notification_id,omitempty"`
	RenewalScheduleId string `json:"renewal_schedule_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	SentAt time.Time `json:"sent_at,omitempty"`
	Status interface{} `json:"status"`
}
