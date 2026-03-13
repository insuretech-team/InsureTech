package models


// UserRolesListingResponse represents a user_roles_listing_response
type UserRolesListingResponse struct {
	UserRoles []*UserRole `json:"user_roles,omitempty"`
	Roles []*Role `json:"roles,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	Error *Error `json:"error,omitempty"`
}
