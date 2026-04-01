package models


// AccessCheckResult represents a access_check_result
type AccessCheckResult struct {
	Action string `json:"action,omitempty"`
	Allowed bool `json:"allowed,omitempty"`
	Object string `json:"object,omitempty"`
	Reason string `json:"reason,omitempty"`
}
