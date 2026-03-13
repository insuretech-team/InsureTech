package models


// SessionsListingRequest represents a sessions_listing_request
type SessionsListingRequest struct {
	DeviceType string `json:"device_type,omitempty"`
	UserId string `json:"user_id"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	ActiveOnly bool `json:"active_only,omitempty"`
}
