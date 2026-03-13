package models


// APIKeysListingResponse represents a api_keys_listing_response
type APIKeysListingResponse struct {
	Error *Error `json:"error,omitempty"`
	Keys []*APIKeySummary `json:"keys,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}
