package models


// VoiceSessionRetrievalResponse represents a voice_session_retrieval_response
type VoiceSessionRetrievalResponse struct {
	Commands []*VoiceCommand `json:"commands,omitempty"`
	VoiceSession *VoiceSession `json:"voice_session,omitempty"`
}
