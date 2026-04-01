package models


// KYCDetection represents a kyc_detection
type KYCDetection struct {
	Detected bool `json:"detected,omitempty"`
	Height int `json:"height,omitempty"`
	Width int `json:"width,omitempty"`
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}
