package models


// PreferencesUpdateRequest represents a preferences_update_request
type PreferencesUpdateRequest struct {
	UserId string `json:"user_id"`
	Preferences *NotificationPreference `json:"preferences,omitempty"`
}
