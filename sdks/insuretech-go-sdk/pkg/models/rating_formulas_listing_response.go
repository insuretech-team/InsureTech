package models


// RatingFormulasListingResponse represents a rating_formulas_listing_response
type RatingFormulasListingResponse struct {
	Formulas []*RatingFormula `json:"formulas,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
