package models


// AccessCheckResult represents a access_check_result
type AccessCheckResult struct {
	Object string `json:"object,omitempty"`
	Action string `json:"action,omitempty"`
	Allowed bool `json:"allowed,omitempty"`
	Reason string `json:"reason,omitempty"`
}
