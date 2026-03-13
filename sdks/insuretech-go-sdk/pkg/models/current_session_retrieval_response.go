package models


// CurrentSessionRetrievalResponse represents a current_session_retrieval_response
type CurrentSessionRetrievalResponse struct {
	Session *Session `json:"session,omitempty"`
	Error *Error `json:"error,omitempty"`
	UserType string `json:"user_type,omitempty"`
}
