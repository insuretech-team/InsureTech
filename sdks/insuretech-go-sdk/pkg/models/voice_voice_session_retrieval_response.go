package models


// VoiceVoiceSessionRetrievalResponse represents a voice_voice_session_retrieval_response
type VoiceVoiceSessionRetrievalResponse struct {
	Commands []*VoiceCommand `json:"commands,omitempty"`
	VoiceSession *VoiceSession `json:"voice_session,omitempty"`
}
