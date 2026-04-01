package models


// KYCEyeState represents a kyc_eye_state
type KYCEyeState struct {
	IsBlinking bool `json:"is_blinking,omitempty"`
	LeftOpenness float64 `json:"left_openness,omitempty"`
	RightOpenness float64 `json:"right_openness,omitempty"`
}
