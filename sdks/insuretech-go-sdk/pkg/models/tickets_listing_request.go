package models


// TicketsListingRequest represents a tickets_listing_request
type TicketsListingRequest struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
}
