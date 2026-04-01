package models


// LeadsListingResponse represents a leads_listing_response
type LeadsListingResponse struct {
	Leads []*Lead `json:"leads,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
