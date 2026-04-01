package models


// CheckAccessRequest represents a check_access_request
type CheckAccessRequest struct {
	Action string `json:"action"`
	Context *AccessContext `json:"context,omitempty"`
	Domain string `json:"domain,omitempty"`
	Object string `json:"object,omitempty"`
	UserId string `json:"user_id"`
}
