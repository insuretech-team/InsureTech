package models

import (
	"time"
)

// IoTDevice represents a iot_device
type IoTDevice struct {
	CreatedAt time.Time `json:"created_at"`
	DeviceId string `json:"device_id"`
	DeviceSerial string `json:"device_serial"`
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
	Manufacturer string `json:"manufacturer"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Model string `json:"model"`
	OwnerId string `json:"owner_id"`
	PolicyId string `json:"policy_id,omitempty"`
	RegisteredAt time.Time `json:"registered_at"`
	Status interface{} `json:"status"`
	Type *IoTDeviceType `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}
