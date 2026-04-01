package models

import (
	"time"
)

// PushDeviceToken represents a push_device_token
type PushDeviceToken struct {
	AppId string `json:"app_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	DeviceId string `json:"device_id,omitempty"`
	DeviceToken string `json:"device_token"`
	IsActive bool `json:"is_active"`
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
	Platform string `json:"platform"`
	Provider string `json:"provider"`
	TokenId string `json:"token_id"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId string `json:"user_id"`
}
