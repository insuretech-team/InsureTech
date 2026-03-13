package models


// CheckAccessResponse represents a check_access_response
type CheckAccessResponse struct {
	Effect *PolicyEffect `json:"effect,omitempty"`
	MatchedRule string `json:"matched_rule,omitempty"`
	Reason string `json:"reason,omitempty"`
	Error *Error `json:"error,omitempty"`
	Allowed bool `json:"allowed,omitempty"`
}
