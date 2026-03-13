package models


// OTPSendingRequest represents a otp_sending_request
type OTPSendingRequest struct {
	Type string `json:"type"`
	Channel string `json:"channel,omitempty"`
	UseMasking bool `json:"use_masking,omitempty"`
	Recipient string `json:"recipient,omitempty"`
}
