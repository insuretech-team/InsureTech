package models


// VehiclesListingRequest represents a vehicles_listing_request
type VehiclesListingRequest struct {
	Filter string `json:"filter,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	OnlyActive bool `json:"only_active,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	VehicleType *VehicleType `json:"vehicle_type"`
	Year int `json:"year,omitempty"`
}
