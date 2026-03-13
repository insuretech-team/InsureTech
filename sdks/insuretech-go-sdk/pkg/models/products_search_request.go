package models


// ProductsSearchRequest represents a products_search_request
type ProductsSearchRequest struct {
	Query string `json:"query"`
	Category *ProductCategory `json:"category,omitempty"`
	MinPremium *Money `json:"min_premium,omitempty"`
	MaxPremium *Money `json:"max_premium,omitempty"`
}
