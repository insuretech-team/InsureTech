package models


// UserProfileRetrievalResponse represents a user_profile_retrieval_response
type UserProfileRetrievalResponse struct {
	Profile *UserProfile `json:"profile,omitempty"`
	Error *Error `json:"error,omitempty"`
}
