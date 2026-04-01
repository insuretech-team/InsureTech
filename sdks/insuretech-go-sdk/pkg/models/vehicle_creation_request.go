package models


// VehicleCreationRequest represents a vehicle_creation_request
type VehicleCreationRequest struct {
	ImageUri string `json:"image_uri,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Model string `json:"model,omitempty"`
	Price string `json:"price,omitempty"`
	Type *VehicleType `json:"type"`
	Year int `json:"year,omitempty"`
}
