package models


// RolesListingRequest represents a roles_listing_request
type RolesListingRequest struct {
	Portal *Portal `json:"portal"`
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
