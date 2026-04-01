package models


// ProductsSearchRequest represents a products_search_request
type ProductsSearchRequest struct {
	Category *ProductCategory `json:"category,omitempty"`
	MaxPremium *Money `json:"max_premium,omitempty"`
	MinPremium *Money `json:"min_premium,omitempty"`
	Query string `json:"query"`
}
