package models

import (
	"time"
)

// AccessDecisionAuditsListingRequest represents a access_decision_audits_listing_request
type AccessDecisionAuditsListingRequest struct {
	Decision *PolicyEffect `json:"decision,omitempty"`
	Domain string `json:"domain,omitempty"`
	From time.Time `json:"from,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	To time.Time `json:"to,omitempty"`
	UserId string `json:"user_id"`
}
