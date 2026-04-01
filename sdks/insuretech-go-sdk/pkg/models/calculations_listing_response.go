package models


// CalculationsListingResponse represents a calculations_listing_response
type CalculationsListingResponse struct {
	Calculations []*ActuarialCalculation `json:"calculations,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
