package models

import (
	"time"
)

// NotificationFailedEvent represents a notification_failed_event
type NotificationFailedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NotificationId string `json:"notification_id,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
