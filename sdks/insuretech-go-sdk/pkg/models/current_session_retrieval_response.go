package models


// CurrentSessionRetrievalResponse represents a current_session_retrieval_response
type CurrentSessionRetrievalResponse struct {
	PasswordChangeRequired bool `json:"password_change_required,omitempty"`
	Session *Session `json:"session,omitempty"`
	UserType string `json:"user_type,omitempty"`
}
