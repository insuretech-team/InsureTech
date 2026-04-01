package models


// AccessDecisionAuditsListingResponse represents a access_decision_audits_listing_response
type AccessDecisionAuditsListingResponse struct {
	Audits []*AccessDecisionAudit `json:"audits,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
