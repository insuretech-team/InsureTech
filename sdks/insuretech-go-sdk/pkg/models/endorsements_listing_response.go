package models


// EndorsementsListingResponse represents a endorsements_listing_response
type EndorsementsListingResponse struct {
	Endorsements []*Endorsement `json:"endorsements,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
