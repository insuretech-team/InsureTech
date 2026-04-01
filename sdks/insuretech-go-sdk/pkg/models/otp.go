package models

import (
	"time"
)

// OTP represents a otp
type OTP struct {
	Attempts int `json:"attempts,omitempty"`
	Carrier string `json:"carrier,omitempty"`
	Channel string `json:"channel,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	DlrErrorCode string `json:"dlr_error_code,omitempty"`
	DlrReceivedAt time.Time `json:"dlr_received_at,omitempty"`
	DlrStatus string `json:"dlr_status,omitempty"`
	DlrUpdatedAt time.Time `json:"dlr_updated_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	OtpHash string `json:"otp_hash,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	ProviderMessageId string `json:"provider_message_id,omitempty"`
	Purpose string `json:"purpose,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Verified bool `json:"verified,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
}
