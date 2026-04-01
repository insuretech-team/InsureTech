package models


// VehicleUpdateRequest represents a vehicle_update_request
type VehicleUpdateRequest struct {
	ImageUri string `json:"image_uri,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Model string `json:"model,omitempty"`
	Price string `json:"price,omitempty"`
	Type *VehicleType `json:"type"`
	VehicleId string `json:"vehicle_id"`
	Year int `json:"year,omitempty"`
}
