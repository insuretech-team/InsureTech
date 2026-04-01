package models


// VehicleDeletionRequest represents a vehicle_deletion_request
type VehicleDeletionRequest struct {
	Permanent bool `json:"permanent,omitempty"`
	VehicleId string `json:"vehicle_id"`
}
