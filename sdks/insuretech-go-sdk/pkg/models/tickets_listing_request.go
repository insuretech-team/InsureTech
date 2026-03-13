package models


// TicketsListingRequest represents a tickets_listing_request
type TicketsListingRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Status string `json:"status,omitempty"`
	AssignedTo string `json:"assigned_to,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
