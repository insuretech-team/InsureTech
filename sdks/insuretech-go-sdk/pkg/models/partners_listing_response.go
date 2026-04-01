package models


// PartnersListingResponse represents a partners_listing_response
type PartnersListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Partners []*Partner `json:"partners,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
