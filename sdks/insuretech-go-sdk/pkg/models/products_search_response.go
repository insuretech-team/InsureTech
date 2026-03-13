package models


// ProductsSearchResponse represents a products_search_response
type ProductsSearchResponse struct {
	Products []*Product `json:"products,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
