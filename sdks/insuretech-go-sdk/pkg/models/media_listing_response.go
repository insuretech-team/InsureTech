package models


// MediaListingResponse represents a media_listing_response
type MediaListingResponse struct {
	MediaFiles []*MediaFile `json:"media_files,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
