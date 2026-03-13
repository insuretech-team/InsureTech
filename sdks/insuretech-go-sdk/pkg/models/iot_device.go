package models

import (
	"time"
)

// IoTDevice represents a iot_device
type IoTDevice struct {
	Type *IotDeviceType `json:"type"`
	Model string `json:"model"`
	PolicyId string `json:"policy_id,omitempty"`
	OwnerId string `json:"owner_id"`
	RegisteredAt time.Time `json:"registered_at"`
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	DeviceSerial string `json:"device_serial"`
	Manufacturer string `json:"manufacturer"`
	Status interface{} `json:"status"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	DeviceId string `json:"device_id"`
}
