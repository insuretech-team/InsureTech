package models


// PreferencesUpdateResponse represents a preferences_update_response
type PreferencesUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
