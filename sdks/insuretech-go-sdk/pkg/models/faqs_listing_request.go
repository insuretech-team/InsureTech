package models


// FAQsListingRequest represents a faqs_listing_request
type FAQsListingRequest struct {
	Category string `json:"category"`
	Language string `json:"language,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
