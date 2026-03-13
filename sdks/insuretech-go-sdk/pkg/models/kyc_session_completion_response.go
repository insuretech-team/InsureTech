package models


// KYCSessionCompletionResponse represents a kyc_session_completion_response
type KYCSessionCompletionResponse struct {
	KycId string `json:"kyc_id,omitempty"`
	Status string `json:"status,omitempty"`
	Success bool `json:"success,omitempty"`
	LivenessConfidence float64 `json:"liveness_confidence,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
