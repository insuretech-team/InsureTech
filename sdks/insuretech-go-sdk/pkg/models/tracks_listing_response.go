package models


// TracksListingResponse represents a tracks_listing_response
type TracksListingResponse struct {
	TotalCount int `json:"total_count,omitempty"`
	Tracks []*Track `json:"tracks,omitempty"`
}
