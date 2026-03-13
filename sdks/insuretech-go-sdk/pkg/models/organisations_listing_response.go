package models


// OrganisationsListingResponse represents a organisations_listing_response
type OrganisationsListingResponse struct {
	Organisations []*Organisation `json:"organisations,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
