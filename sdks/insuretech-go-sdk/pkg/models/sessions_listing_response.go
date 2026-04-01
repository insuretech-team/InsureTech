package models


// SessionsListingResponse represents a sessions_listing_response
type SessionsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Sessions []*Session `json:"sessions,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
