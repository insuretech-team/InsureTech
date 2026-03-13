package models


// UserProfileUpdateResponse represents a user_profile_update_response
type UserProfileUpdateResponse struct {
	Profile *UserProfile `json:"profile,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
