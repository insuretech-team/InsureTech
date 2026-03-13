package models


// RolesListingResponse represents a roles_listing_response
type RolesListingResponse struct {
	Error *Error `json:"error,omitempty"`
	Roles []*Role `json:"roles,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
