package models


// PendingVerificationsListingResponse represents a pending_verifications_listing_response
type PendingVerificationsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Verifications []*KYCVerification `json:"verifications,omitempty"`
}
