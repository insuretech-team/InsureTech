package models


// NotificationPreferencesRetrievalResponse represents a notification_preferences_retrieval_response
type NotificationPreferencesRetrievalResponse struct {
	NotificationPreference string `json:"notification_preference,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
