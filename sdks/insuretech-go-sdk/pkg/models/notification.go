package models

import (
	"time"
)

// Notification represents a notification
type Notification struct {
	Channel *NotificationChannel `json:"channel"`
	CreatedAt time.Time `json:"created_at"`
	DeliveredAt time.Time `json:"delivered_at,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Message string `json:"message"`
	NotificationId string `json:"notification_id"`
	Priority interface{} `json:"priority"`
	ReadAt time.Time `json:"read_at,omitempty"`
	RecipientId string `json:"recipient_id"`
	RetryCount int `json:"retry_count"`
	ScheduledAt time.Time `json:"scheduled_at,omitempty"`
	SentAt time.Time `json:"sent_at,omitempty"`
	Status interface{} `json:"status"`
	Subject string `json:"subject,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Type *NotificationType `json:"type"`
}
