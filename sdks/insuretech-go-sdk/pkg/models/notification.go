package models

import (
	"time"
)

// Notification represents a notification
type Notification struct {
	Type *NotificationType `json:"type"`
	Channel *NotificationChannel `json:"channel"`
	SentAt time.Time `json:"sent_at,omitempty"`
	DeliveredAt time.Time `json:"delivered_at,omitempty"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message"`
	Priority interface{} `json:"priority"`
	Status interface{} `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	NotificationId string `json:"notification_id"`
	ScheduledAt time.Time `json:"scheduled_at,omitempty"`
	ReadAt time.Time `json:"read_at,omitempty"`
	RetryCount int `json:"retry_count"`
	ErrorMessage string `json:"error_message,omitempty"`
	RecipientId string `json:"recipient_id"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
}
