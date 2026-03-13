package models

import (
	"time"
)

// NotificationDeliveredEvent represents a notification_delivered_event
type NotificationDeliveredEvent struct {
	EventId string `json:"event_id,omitempty"`
	NotificationId string `json:"notification_id,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
