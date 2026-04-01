package models


// OrganisationsListingResponse represents a organisations_listing_response
type OrganisationsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Organisations []*Organisation `json:"organisations,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
