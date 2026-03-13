package models


// NotificationPreferencesUpdateResponse represents a notification_preferences_update_response
type NotificationPreferencesUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
