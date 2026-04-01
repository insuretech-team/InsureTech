package models


// KYCStepResult represents a kyc_step_result
type KYCStepResult struct {
	ChallengeType string `json:"challenge_type,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	ElapsedMs int `json:"elapsed_ms,omitempty"`
	FramesProcessed int `json:"frames_processed,omitempty"`
	State string `json:"state,omitempty"`
}
