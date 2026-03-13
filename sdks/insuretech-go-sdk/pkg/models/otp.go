package models

import (
	"time"
)

// OTP represents a otp
type OTP struct {
	SenderId string `json:"sender_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	Attempts int `json:"attempts,omitempty"`
	DlrStatus string `json:"dlr_status,omitempty"`
	Carrier string `json:"carrier,omitempty"`
	Channel string `json:"channel,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Verified bool `json:"verified,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	ProviderMessageId string `json:"provider_message_id,omitempty"`
	DlrReceivedAt time.Time `json:"dlr_received_at,omitempty"`
	DlrErrorCode string `json:"dlr_error_code,omitempty"`
	OtpHash string `json:"otp_hash,omitempty"`
	Purpose string `json:"purpose,omitempty"`
	DlrUpdatedAt time.Time `json:"dlr_updated_at,omitempty"`
}
