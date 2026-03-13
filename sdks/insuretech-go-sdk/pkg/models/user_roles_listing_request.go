package models


// UserRolesListingRequest represents a user_roles_listing_request
type UserRolesListingRequest struct {
	UserId string `json:"user_id"`
	Domain string `json:"domain,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
