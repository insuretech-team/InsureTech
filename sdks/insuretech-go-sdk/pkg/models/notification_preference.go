package models

import (
	"time"
)

// NotificationPreference represents a notification_preference
type NotificationPreference struct {
	MarketingOptIn bool `json:"marketing_opt_in,omitempty"`
	TransactionalOptIn bool `json:"transactional_opt_in,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserId string `json:"user_id,omitempty"`
	ChannelPreferences []*ChannelPreference `json:"channel_preferences,omitempty"`
}
