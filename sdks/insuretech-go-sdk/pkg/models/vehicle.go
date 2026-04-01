package models

import (
	"time"
)

// Vehicle represents a vehicle
type Vehicle struct {
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	ImageUri string `json:"image_uri,omitempty"`
	IsActive bool `json:"is_active"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Model string `json:"model"`
	Price string `json:"price"`
	Type *VehicleType `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	VehicleId string `json:"vehicle_id"`
	Year int `json:"year,omitempty"`
}
