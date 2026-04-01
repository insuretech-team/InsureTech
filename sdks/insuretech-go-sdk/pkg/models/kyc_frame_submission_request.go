package models


// KYCFrameSubmissionRequest represents a kyc_frame_submission_request
type KYCFrameSubmissionRequest struct {
	FrameSequence int `json:"frame_sequence,omitempty"`
	ImageData string `json:"image_data,omitempty"`
	SessionId string `json:"session_id"`
	UserId string `json:"user_id"`
}
