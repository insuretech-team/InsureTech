package models


// TracksListingResponse represents a tracks_listing_response
type TracksListingResponse struct {
	Tracks []*Track `json:"tracks,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
