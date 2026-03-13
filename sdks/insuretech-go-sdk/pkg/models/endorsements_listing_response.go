package models


// EndorsementsListingResponse represents a endorsements_listing_response
type EndorsementsListingResponse struct {
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	Endorsements []*Endorsement `json:"endorsements,omitempty"`
}
