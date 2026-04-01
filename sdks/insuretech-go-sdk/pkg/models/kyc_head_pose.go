package models


// KYCHeadPose represents a kyc_head_pose
type KYCHeadPose struct {
	Pitch float64 `json:"pitch,omitempty"`
	Roll float64 `json:"roll,omitempty"`
	Yaw float64 `json:"yaw,omitempty"`
}
