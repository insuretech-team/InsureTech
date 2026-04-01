package models


// FindPortalUserResponse represents a find_portal_user_response
type FindPortalUserResponse struct {
	User *PortalUserSummary `json:"user,omitempty"`
}
