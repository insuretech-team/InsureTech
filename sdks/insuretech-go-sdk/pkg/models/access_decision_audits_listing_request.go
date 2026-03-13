package models

import (
	"time"
)

// AccessDecisionAuditsListingRequest represents a access_decision_audits_listing_request
type AccessDecisionAuditsListingRequest struct {
	UserId string `json:"user_id"`
	Domain string `json:"domain,omitempty"`
	Decision *PolicyEffect `json:"decision,omitempty"`
	From time.Time `json:"from,omitempty"`
	To time.Time `json:"to,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
