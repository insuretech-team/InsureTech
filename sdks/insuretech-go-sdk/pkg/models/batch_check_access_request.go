package models


// BatchCheckAccessRequest represents a batch_check_access_request
type BatchCheckAccessRequest struct {
	Domain string `json:"domain,omitempty"`
	Checks []*AccessCheckTuple `json:"checks,omitempty"`
	Context *AccessContext `json:"context,omitempty"`
	UserId string `json:"user_id"`
}
