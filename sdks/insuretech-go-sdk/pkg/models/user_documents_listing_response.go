package models


// UserDocumentsListingResponse represents a user_documents_listing_response
type UserDocumentsListingResponse struct {
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	Documents []*UserDocument `json:"documents,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}
