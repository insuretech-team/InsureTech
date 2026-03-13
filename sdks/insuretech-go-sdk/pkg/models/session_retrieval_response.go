package models


// SessionRetrievalResponse represents a session_retrieval_response
type SessionRetrievalResponse struct {
	Session *Session `json:"session,omitempty"`
	Error *Error `json:"error,omitempty"`
}
