package models


// LossRatioCalculationsListingResponse represents a loss_ratio_calculations_listing_response
type LossRatioCalculationsListingResponse struct {
	LossRatios []*LossRatioCalculation `json:"loss_ratios,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
