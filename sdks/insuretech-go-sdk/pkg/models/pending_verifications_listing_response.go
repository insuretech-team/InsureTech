package models


// PendingVerificationsListingResponse represents a pending_verifications_listing_response
type PendingVerificationsListingResponse struct {
	Verifications []*KYCVerification `json:"verifications,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
