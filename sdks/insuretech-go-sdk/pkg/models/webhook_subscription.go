package models

import (
	"time"
)

// WebhookSubscription represents a webhook_subscription
type WebhookSubscription struct {
	Channels []string `json:"channels,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	EventTypes []string `json:"event_types,omitempty"`
	IsActive bool `json:"is_active"`
	MaxAttempts int `json:"max_attempts"`
	Secret string `json:"secret"`
	SubscriberName string `json:"subscriber_name"`
	SubscriptionId string `json:"subscription_id"`
	TargetUrl string `json:"target_url"`
	TimeoutSeconds int `json:"timeout_seconds"`
	TopicGroups []string `json:"topic_groups,omitempty"`
	Topics []string `json:"topics,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
