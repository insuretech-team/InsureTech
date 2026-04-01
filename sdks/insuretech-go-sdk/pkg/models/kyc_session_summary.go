package models


// KYCSessionSummary represents a kyc_session_summary
type KYCSessionSummary struct {
	CompletedSteps int `json:"completed_steps,omitempty"`
	ElapsedMs int `json:"elapsed_ms,omitempty"`
	FailedSteps int `json:"failed_steps,omitempty"`
	StepResults []*KYCStepResult `json:"step_results,omitempty"`
	TotalFramesProcessed int `json:"total_frames_processed,omitempty"`
	TotalSteps int `json:"total_steps,omitempty"`
}
