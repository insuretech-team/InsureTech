package models


// CheckAccessResponse represents a check_access_response
type CheckAccessResponse struct {
	Allowed bool `json:"allowed,omitempty"`
	Effect *PolicyEffect `json:"effect,omitempty"`
	MatchedRule string `json:"matched_rule,omitempty"`
	Reason string `json:"reason,omitempty"`
}
