package models


// VoiceSessionRetrievalResponse represents a voice_session_retrieval_response
type VoiceSessionRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	VoiceSession *VoiceSession `json:"voice_session,omitempty"`
	Commands []*VoiceCommand `json:"commands,omitempty"`
}
