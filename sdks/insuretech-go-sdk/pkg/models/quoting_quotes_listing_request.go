package models


// QuotingQuotesListingRequest represents a quoting_quotes_listing_request
type QuotingQuotesListingRequest struct {
	AgentId string `json:"agent_id"`
	CustomerId string `json:"customer_id"`
	Filter string `json:"filter,omitempty"`
	OrderBy string `json:"order_by,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	ProductId string `json:"product_id"`
	Status *QuotingQuoteStatus `json:"status,omitempty"`
}
