package models


// InsuranceProposalsListingResponse represents a insurance_proposals_listing_response
type InsuranceProposalsListingResponse struct {
	Proposals []*InsuranceProposal `json:"proposals,omitempty"`
	Total int `json:"total,omitempty"`
}
