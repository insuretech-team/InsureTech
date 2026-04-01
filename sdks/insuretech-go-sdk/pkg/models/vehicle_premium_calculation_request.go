package models


// VehiclePremiumCalculationRequest represents a vehicle_premium_calculation_request
type VehiclePremiumCalculationRequest struct {
	AccidentalCover bool `json:"accidental_cover,omitempty"`
	RegistrationId string `json:"registration_id"`
	RegistrationNumber string `json:"registration_number,omitempty"`
}
