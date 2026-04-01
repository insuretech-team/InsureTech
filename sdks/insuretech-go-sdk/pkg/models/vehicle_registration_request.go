package models


// VehicleRegistrationRequest represents a vehicle_registration_request
type VehicleRegistrationRequest struct {
	AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
	OwnerId string `json:"owner_id"`
	RegistrationNumber string `json:"registration_number,omitempty"`
	RegistrationYear int `json:"registration_year,omitempty"`
	VehicleId string `json:"vehicle_id"`
}
