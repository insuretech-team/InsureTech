package models


// PendingVerificationsListingRequest represents a pending_verifications_listing_request
type PendingVerificationsListingRequest struct {
	PageSize int `json:"page_size"`
	PageToken string `json:"page_token,omitempty"`
}
