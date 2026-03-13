package models


// NotificationPreferencesUpdateRequest represents a notification_preferences_update_request
type NotificationPreferencesUpdateRequest struct {
	PreferredLanguage string `json:"preferred_language,omitempty"`
	UserId string `json:"user_id"`
	NotificationPreference string `json:"notification_preference,omitempty"`
}
