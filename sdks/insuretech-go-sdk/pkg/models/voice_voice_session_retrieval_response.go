package models


// VoiceVoiceSessionRetrievalResponse represents a voice_voice_session_retrieval_response
type VoiceVoiceSessionRetrievalResponse struct {
	VoiceSession *VoiceSession `json:"voice_session,omitempty"`
	Commands []*VoiceCommand `json:"commands,omitempty"`
	Error *Error `json:"error,omitempty"`
}
