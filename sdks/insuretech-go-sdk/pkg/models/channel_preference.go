package models


// ChannelPreference represents a channel_preference
type ChannelPreference struct {
	Channel *NotificationChannel `json:"channel,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
	ExcludedTypes []*NotificationType `json:"excluded_types,omitempty"`
}
