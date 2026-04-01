package models


// VehiclesListingResponse represents a vehicles_listing_response
type VehiclesListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Vehicles []*Vehicle `json:"vehicles,omitempty"`
}
