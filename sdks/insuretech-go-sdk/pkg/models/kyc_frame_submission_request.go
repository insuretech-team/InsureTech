package models


// KYCFrameSubmissionRequest represents a kyc_frame_submission_request
type KYCFrameSubmissionRequest struct {
	FrameSequence int `json:"frame_sequence,omitempty"`
	UserId string `json:"user_id"`
	SessionId string `json:"session_id"`
	ImageData string `json:"image_data,omitempty"`
}
