package models

import (
	"time"
)

// NotificationSentEvent represents a notification_sent_event
type NotificationSentEvent struct {
	NotificationId string `json:"notification_id,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	Channel string `json:"channel,omitempty"`
	Type string `json:"type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
