package models


// PolicyUpdateResponse represents a policy_update_response
type PolicyUpdateResponse struct {
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
