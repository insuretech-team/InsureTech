package models


// KYCSessionCompletionResponse represents a kyc_session_completion_response
type KYCSessionCompletionResponse struct {
	CompletedAt string `json:"completed_at,omitempty"`
	IdentityMatch bool `json:"identity_match,omitempty"`
	KycId string `json:"kyc_id,omitempty"`
	LivenessConfidence float64 `json:"liveness_confidence,omitempty"`
	MatchScore float64 `json:"match_score,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	ProviderReference string `json:"provider_reference,omitempty"`
	SessionState string `json:"session_state,omitempty"`
	Status string `json:"status,omitempty"`
	Success bool `json:"success,omitempty"`
	Summary *KYCSessionSummary `json:"summary,omitempty"`
}
