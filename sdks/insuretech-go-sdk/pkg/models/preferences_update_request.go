package models


// PreferencesUpdateRequest represents a preferences_update_request
type PreferencesUpdateRequest struct {
	Preferences *NotificationPreference `json:"preferences,omitempty"`
	UserId string `json:"user_id"`
}
