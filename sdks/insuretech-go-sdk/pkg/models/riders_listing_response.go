package models


// RidersListingResponse represents a riders_listing_response
type RidersListingResponse struct {
	Riders []*ProductsRider `json:"riders,omitempty"`
}
