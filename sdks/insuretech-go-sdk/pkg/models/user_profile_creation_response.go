package models


// UserProfileCreationResponse represents a user_profile_creation_response
type UserProfileCreationResponse struct {
	Profile *UserProfile `json:"profile,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
