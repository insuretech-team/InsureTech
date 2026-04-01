package models


// VoiceSampleSubmissionRequest represents a voice_sample_submission_request
type VoiceSampleSubmissionRequest struct {
	AudioUrl string `json:"audio_url,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	SessionId string `json:"session_id"`
	Transcript string `json:"transcript,omitempty"`
}
