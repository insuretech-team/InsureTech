package models

import (
	"time"
)

// WebhookDeliveryAttempt represents a webhook_delivery_attempt
type WebhookDeliveryAttempt struct {
	AttemptId string `json:"attempt_id"`
	CreatedAt time.Time `json:"created_at"`
	ErrorMessage string `json:"error_message,omitempty"`
	LastAttemptedAt time.Time `json:"last_attempted_at,omitempty"`
	LifecycleEvent string `json:"lifecycle_event"`
	NotificationId string `json:"notification_id,omitempty"`
	Payload string `json:"payload"`
	ResponseBody string `json:"response_body,omitempty"`
	ResponseStatus int `json:"response_status,omitempty"`
	RetryCount int `json:"retry_count"`
	ScheduledAt time.Time `json:"scheduled_at,omitempty"`
	SourceTopic string `json:"source_topic,omitempty"`
	Status string `json:"status"`
	SubscriptionId string `json:"subscription_id"`
	UpdatedAt time.Time `json:"updated_at"`
}
