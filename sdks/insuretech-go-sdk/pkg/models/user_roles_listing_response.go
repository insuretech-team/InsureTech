package models


// UserRolesListingResponse represents a user_roles_listing_response
type UserRolesListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Roles []*Role `json:"roles,omitempty"`
	UserRoles []*UserRole `json:"user_roles,omitempty"`
}
