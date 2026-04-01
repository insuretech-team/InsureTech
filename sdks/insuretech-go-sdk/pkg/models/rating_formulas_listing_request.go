package models


// RatingFormulasListingRequest represents a rating_formulas_listing_request
type RatingFormulasListingRequest struct {
	Category *FormulaCategory `json:"category,omitempty"`
	InsuranceType string `json:"insurance_type"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	SearchQuery string `json:"search_query,omitempty"`
	Status *FormulaStatus `json:"status,omitempty"`
}
