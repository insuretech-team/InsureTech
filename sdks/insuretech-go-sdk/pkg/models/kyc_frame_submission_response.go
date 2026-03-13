package models


// KYCFrameSubmissionResponse represents a kyc_frame_submission_response
type KYCFrameSubmissionResponse struct {
	TotalSteps int `json:"total_steps,omitempty"`
	LivenessConfidence float64 `json:"liveness_confidence,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	Accepted bool `json:"accepted,omitempty"`
	Guidance string `json:"guidance,omitempty"`
	CurrentStep string `json:"current_step,omitempty"`
	CompletedSteps int `json:"completed_steps,omitempty"`
}
