package models


// LifeQuotesListingRequest represents a life_quotes_listing_request
type LifeQuotesListingRequest struct {
	CustomerId string `json:"customer_id"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	ProductId string `json:"product_id"`
	Status *LifeQuoteStatus `json:"status,omitempty"`
}
