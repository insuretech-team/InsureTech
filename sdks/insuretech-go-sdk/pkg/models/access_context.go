package models


// AccessContext represents a access_context
type AccessContext struct {
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	KycVerified bool `json:"kyc_verified,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	TokenId string `json:"token_id,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
