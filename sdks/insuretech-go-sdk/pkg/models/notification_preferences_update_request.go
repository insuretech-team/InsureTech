package models


// NotificationPreferencesUpdateRequest represents a notification_preferences_update_request
type NotificationPreferencesUpdateRequest struct {
	NotificationPreference string `json:"notification_preference,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	UserId string `json:"user_id"`
}
