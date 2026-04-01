package models


// ContactsListingResponse represents a contacts_listing_response
type ContactsListingResponse struct {
	Contacts []*Contact `json:"contacts,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
