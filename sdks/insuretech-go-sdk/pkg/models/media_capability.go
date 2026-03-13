package models


// MediaCapability represents a media_capability
type MediaCapability struct {
	MimeType string `json:"mime_type,omitempty"`
	ClockRate int `json:"clock_rate,omitempty"`
	Channels int `json:"channels,omitempty"`
	SdpFmtpLine string `json:"sdp_fmtp_line,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}
