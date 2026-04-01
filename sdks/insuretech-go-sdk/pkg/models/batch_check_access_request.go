package models


// BatchCheckAccessRequest represents a batch_check_access_request
type BatchCheckAccessRequest struct {
	Checks []*AccessCheckTuple `json:"checks,omitempty"`
	Context *AccessContext `json:"context,omitempty"`
	Domain string `json:"domain,omitempty"`
	UserId string `json:"user_id"`
}
