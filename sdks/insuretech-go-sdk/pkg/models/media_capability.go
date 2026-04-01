package models


// MediaCapability represents a media_capability
type MediaCapability struct {
	Channels int `json:"channels,omitempty"`
	ClockRate int `json:"clock_rate,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	SdpFmtpLine string `json:"sdp_fmtp_line,omitempty"`
}
