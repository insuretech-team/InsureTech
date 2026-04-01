package models


// UserRolesListingRequest represents a user_roles_listing_request
type UserRolesListingRequest struct {
	Domain string `json:"domain,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	UserId string `json:"user_id"`
}
