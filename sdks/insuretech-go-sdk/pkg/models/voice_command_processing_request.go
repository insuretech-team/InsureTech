package models


// VoiceCommandProcessingRequest represents a voice_command_processing_request
type VoiceCommandProcessingRequest struct {
	AudioUrl string `json:"audio_url,omitempty"`
	CommandText string `json:"command_text,omitempty"`
	Language string `json:"language,omitempty"`
	VoiceSessionId string `json:"voice_session_id"`
}
