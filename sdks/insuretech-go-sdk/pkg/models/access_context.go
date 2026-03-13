package models


// AccessContext represents a access_context
type AccessContext struct {
	IpAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	TokenId string `json:"token_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	KycVerified bool `json:"kyc_verified,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}
