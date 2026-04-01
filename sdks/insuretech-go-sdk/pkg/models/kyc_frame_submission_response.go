package models


// KYCFrameSubmissionResponse represents a kyc_frame_submission_response
type KYCFrameSubmissionResponse struct {
	Accepted bool `json:"accepted,omitempty"`
	CompletedSteps int `json:"completed_steps,omitempty"`
	CurrentStep string `json:"current_step,omitempty"`
	CurrentStepDetail *KYCStep `json:"current_step_detail,omitempty"`
	Detection *KYCDetection `json:"detection,omitempty"`
	EyeContoursJson string `json:"eye_contours_json,omitempty"`
	EyeState *KYCEyeState `json:"eye_state,omitempty"`
	Guidance string `json:"guidance,omitempty"`
	GuidanceMessages []string `json:"guidance_messages,omitempty"`
	HeadPose *KYCHeadPose `json:"head_pose,omitempty"`
	LivenessConfidence float64 `json:"liveness_confidence,omitempty"`
	OverallProgress float64 `json:"overall_progress,omitempty"`
	SessionState string `json:"session_state,omitempty"`
	StepCompleted bool `json:"step_completed,omitempty"`
	StepProgress float64 `json:"step_progress,omitempty"`
	TotalSteps int `json:"total_steps,omitempty"`
}
