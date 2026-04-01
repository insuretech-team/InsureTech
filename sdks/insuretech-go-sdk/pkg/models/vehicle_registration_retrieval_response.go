package models


// VehicleRegistrationRetrievalResponse represents a vehicle_registration_retrieval_response
type VehicleRegistrationRetrievalResponse struct {
	Registration *VehicleRegistration `json:"registration,omitempty"`
	Vehicle *Vehicle `json:"vehicle,omitempty"`
}
