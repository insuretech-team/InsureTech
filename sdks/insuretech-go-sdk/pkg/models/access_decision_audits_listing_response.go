package models


// AccessDecisionAuditsListingResponse represents a access_decision_audits_listing_response
type AccessDecisionAuditsListingResponse struct {
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	Audits []*AccessDecisionAudit `json:"audits,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}
