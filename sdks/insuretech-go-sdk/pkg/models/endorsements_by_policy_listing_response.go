package models


// EndorsementsByPolicyListingResponse represents a endorsements_by_policy_listing_response
type EndorsementsByPolicyListingResponse struct {
	Endorsements []*Endorsement `json:"endorsements,omitempty"`
}
