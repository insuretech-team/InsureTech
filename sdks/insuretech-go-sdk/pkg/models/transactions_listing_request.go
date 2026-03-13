package models


// TransactionsListingRequest represents a transactions_listing_request
type TransactionsListingRequest struct {
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Provider string `json:"provider,omitempty"`
	Status string `json:"status,omitempty"`
	PaymentId string `json:"payment_id"`
	StartDate string `json:"start_date,omitempty"`
}
