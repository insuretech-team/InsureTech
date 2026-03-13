package models


// VoiceSampleSubmissionResponse represents a voice_sample_submission_response
type VoiceSampleSubmissionResponse struct {
	Verified bool `json:"verified,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
