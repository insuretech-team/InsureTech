package models


// VoiceCommandProcessingResponse represents a voice_command_processing_response
type VoiceCommandProcessingResponse struct {
	CommandType string `json:"command_type,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	ResponseAudioUrl string `json:"response_audio_url,omitempty"`
	ResponseText string `json:"response_text,omitempty"`
	VoiceCommandId string `json:"voice_command_id,omitempty"`
}
