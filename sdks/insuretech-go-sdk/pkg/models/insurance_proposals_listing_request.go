package models


// InsuranceProposalsListingRequest represents a insurance_proposals_listing_request
type InsuranceProposalsListingRequest struct {
	CustomerId string `json:"customer_id"`
	InsurerId string `json:"insurer_id"`
	OrderId string `json:"order_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status *ProposalStatus `json:"status,omitempty"`
}
