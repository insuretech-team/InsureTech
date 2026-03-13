package models


// UnderwritingDecisionsListingResponse represents a underwriting_decisions_listing_response
type UnderwritingDecisionsListingResponse struct {
	Decisions []*UnderwritingDecision `json:"decisions,omitempty"`
}
