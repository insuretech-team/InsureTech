package models


// CheckAccessRequest represents a check_access_request
type CheckAccessRequest struct {
	UserId string `json:"user_id"`
	Domain string `json:"domain,omitempty"`
	Object string `json:"object,omitempty"`
	Action string `json:"action"`
	Context *AccessContext `json:"context,omitempty"`
}
